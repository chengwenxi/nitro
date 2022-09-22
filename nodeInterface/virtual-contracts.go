// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package nodeInterface

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/mantle"
	"github.com/mantlenetworkio/mantle/mtnode"
	"github.com/mantlenetworkio/mantle/mtos/mtosState"
	"github.com/mantlenetworkio/mantle/mtstate"
	"github.com/mantlenetworkio/mantle/precompiles"
	"github.com/mantlenetworkio/mantle/solgen/go/node_interfacegen"
	"github.com/mantlenetworkio/mantle/solgen/go/precompilesgen"
	"github.com/mantlenetworkio/mantle/util/mtmath"
)

type addr = common.Address
type mech = *vm.EVM
type huge = *big.Int
type hash = common.Hash
type bytes32 = [32]byte
type ctx = *precompiles.Context

type Message = types.Message
type BackendAPI = core.NodeInterfaceBackendAPI
type ExecutionResult = core.ExecutionResult

func init() {
	mtstate.RequireHookedGeth()

	nodeInterfaceImpl := &NodeInterface{Address: types.NodeInterfaceAddress}
	nodeInterfaceMeta := node_interfacegen.NodeInterfaceMetaData
	_, nodeInterface := precompiles.MakePrecompile(nodeInterfaceMeta, nodeInterfaceImpl)

	nodeInterfaceDebugImpl := &NodeInterfaceDebug{Address: types.NodeInterfaceDebugAddress}
	nodeInterfaceDebugMeta := node_interfacegen.NodeInterfaceDebugMetaData
	_, nodeInterfaceDebug := precompiles.MakePrecompile(nodeInterfaceDebugMeta, nodeInterfaceDebugImpl)

	core.InterceptRPCMessage = func(
		msg Message,
		ctx context.Context,
		statedb *state.StateDB,
		header *types.Header,
		backend core.NodeInterfaceBackendAPI,
	) (Message, *ExecutionResult, error) {
		to := msg.To()
		mtosVersion := mtosState.MtOSVersion(statedb) // check MtOS has been installed
		if to != nil && mtosVersion != 0 {
			var precompile precompiles.MtosPrecompile
			var swapMessages bool
			returnMessage := &Message{}
			var address addr

			switch *to {
			case types.NodeInterfaceAddress:
				duplicate := *nodeInterfaceImpl
				duplicate.backend = backend
				duplicate.context = ctx
				duplicate.header = header
				duplicate.sourceMessage = msg
				duplicate.returnMessage.message = returnMessage
				duplicate.returnMessage.changed = &swapMessages
				precompile = nodeInterface.SwapImpl(&duplicate)
				address = types.NodeInterfaceAddress
			case types.NodeInterfaceDebugAddress:
				duplicate := *nodeInterfaceDebugImpl
				duplicate.backend = backend
				duplicate.context = ctx
				duplicate.header = header
				duplicate.sourceMessage = msg
				duplicate.returnMessage.message = returnMessage
				duplicate.returnMessage.changed = &swapMessages
				precompile = nodeInterfaceDebug.SwapImpl(&duplicate)
				address = types.NodeInterfaceDebugAddress
			default:
				return msg, nil, nil
			}

			evm, vmError, err := backend.GetEVM(ctx, msg, statedb, header, &vm.Config{NoBaseFee: true})
			if err != nil {
				return msg, nil, err
			}
			go func() {
				<-ctx.Done()
				evm.Cancel()
			}()
			core.ReadyEVMForL2(evm, msg)

			output, gasLeft, err := precompile.Call(
				msg.Data(), address, address, msg.From(), msg.Value(), false, msg.Gas(), evm,
			)
			if err != nil {
				return msg, nil, err
			}
			if swapMessages {
				return *returnMessage, nil, nil
			}
			res := &ExecutionResult{
				UsedGas:       msg.Gas() - gasLeft,
				Err:           nil,
				ReturnData:    output,
				ScheduledTxes: nil,
			}
			return msg, res, vmError()
		}
		return msg, nil, nil
	}

	core.InterceptRPCGasCap = func(gascap *uint64, msg Message, header *types.Header, statedb *state.StateDB) {
		if *gascap == 0 {
			// It's already unlimited
			return
		}
		mtosVersion := mtosState.MtOSVersion(statedb)
		if mtosVersion == 0 {
			// MtOS hasn't been installed, so use the vanilla gas cap
			return
		}
		state, err := mtosState.OpenSystemMtosState(statedb, nil, true)
		if err != nil {
			log.Error("failed to open MtOS state", "err", err)
			return
		}
		if header.BaseFee.Sign() == 0 {
			// if gas is free or there's no reimbursable poster, the user won't pay for L1 data costs
			return
		}

		posterCost, _ := state.L1PricingState().PosterDataCost(msg, header.Coinbase)
		posterCostInL2Gas := mtmath.BigToUintSaturating(mtmath.BigDiv(posterCost, header.BaseFee))
		*gascap = mtmath.SaturatingUAdd(*gascap, posterCostInL2Gas)
	}

	core.GetMtOSSpeedLimitPerSecond = func(statedb *state.StateDB) (uint64, error) {
		mtosVersion := mtosState.MtOSVersion(statedb)
		if mtosVersion == 0 {
			return 0.0, errors.New("MtOS not installed")
		}
		state, err := mtosState.OpenSystemMtosState(statedb, nil, true)
		if err != nil {
			log.Error("failed to open MtOS state", "err", err)
			return 0.0, err
		}
		pricing := state.L2PricingState()
		speedLimit, err := pricing.SpeedLimitPerSecond()
		if err != nil {
			log.Error("failed to get the speed limit", "err", err)
			return 0.0, err
		}
		return speedLimit, nil
	}

	mtSys, err := precompilesgen.MtSysMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	l2ToL1TxTopic = mtSys.Events["L2ToL1Tx"].ID
	l2ToL1TransactionTopic = mtSys.Events["L2ToL1Transaction"].ID
	merkleTopic = mtSys.Events["SendMerkleUpdate"].ID
}

func mtNodeFromNodeInterfaceBackend(backend BackendAPI) (*mtnode.Node, error) {
	apiBackend, ok := backend.(*mantle.APIBackend)
	if !ok {
		return nil, errors.New("API backend isn't Mantle")
	}
	mtNode, ok := apiBackend.GetMantleNode().(*mtnode.Node)
	if !ok {
		return nil, errors.New("failed to get Mantle Node from backend")
	}
	return mtNode, nil
}
