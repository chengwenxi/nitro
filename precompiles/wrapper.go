// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package precompiles

import (
	"errors"
	"math/big"

	"github.com/mantlenetworkio/mantle/mtos/mtosState"
	"github.com/mantlenetworkio/mantle/mtos/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

// A precompile wrapper for those not allowed in production
type DebugPrecompile struct {
	precompile MtosPrecompile
}

// create a debug-only precompile wrapper
func debugOnly(address addr, impl MtosPrecompile) (addr, MtosPrecompile) {
	return address, &DebugPrecompile{impl}
}

func (wrapper *DebugPrecompile) Call(
	input []byte,
	precompileAddress common.Address,
	actingAsAddress common.Address,
	caller common.Address,
	value *big.Int,
	readOnly bool,
	gasSupplied uint64,
	evm *vm.EVM,
) ([]byte, uint64, error) {

	debugMode := evm.ChainConfig().DebugMode()

	if debugMode {
		con := wrapper.precompile
		return con.Call(input, precompileAddress, actingAsAddress, caller, value, readOnly, gasSupplied, evm)
	} else {
		// take all gas
		return nil, 0, errors.New("Debug precompiles are disabled")
	}
}

func (wrapper *DebugPrecompile) Precompile() Precompile {
	return wrapper.precompile.Precompile()
}

// A precompile wrapper for those only chain owners may use
type OwnerPrecompile struct {
	precompile  MtosPrecompile
	emitSuccess func(mech, bytes4, addr, []byte) error
}

func ownerOnly(address addr, impl MtosPrecompile, emit func(mech, bytes4, addr, []byte) error) (addr, MtosPrecompile) {
	return address, &OwnerPrecompile{
		precompile:  impl,
		emitSuccess: emit,
	}
}

func (wrapper *OwnerPrecompile) Call(
	input []byte,
	precompileAddress common.Address,
	actingAsAddress common.Address,
	caller common.Address,
	value *big.Int,
	readOnly bool,
	gasSupplied uint64,
	evm *vm.EVM,
) ([]byte, uint64, error) {
	con := wrapper.precompile

	burner := &Context{
		gasSupplied: gasSupplied,
		gasLeft:     gasSupplied,
		tracingInfo: util.NewTracingInfo(evm, caller, precompileAddress, util.TracingDuringEVM),
	}
	state, err := mtosState.OpenMtosState(evm.StateDB, burner)
	if err != nil {
		return nil, burner.gasLeft, err
	}

	owners := state.ChainOwners()
	isOwner, err := owners.IsMember(caller)
	if err != nil {
		return nil, burner.gasLeft, err
	}

	if !isOwner {
		return nil, burner.gasLeft, errors.New("unauthorized caller to access-controlled method")
	}

	output, _, err := con.Call(input, precompileAddress, actingAsAddress, caller, value, readOnly, gasSupplied, evm)

	if err != nil {
		return output, gasSupplied, err // we don't deduct gas since we don't want to charge the owner
	}

	// log that the owner operation succeeded
	if err := wrapper.emitSuccess(evm, *(*[4]byte)(input[:4]), caller, input); err != nil {
		log.Error("failed to emit OwnerActs event", "err", err)
	}

	return output, gasSupplied, err // we don't deduct gas since we don't want to charge the owner
}

func (wrapper *OwnerPrecompile) Precompile() Precompile {
	con := wrapper.precompile
	return con.Precompile()
}
