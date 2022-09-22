// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantlenetworkio/mantle/blob/main/LICENSE

package mttest

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mantlenetworkio/mantle/mtos/burn"
	"github.com/mantlenetworkio/mantle/mtos/mtosState"
	"github.com/mantlenetworkio/mantle/mtstate"
	"github.com/mantlenetworkio/mantle/precompiles"
)

const fuzzGas uint64 = 1200000

func FuzzPrecompiles(f *testing.F) {
	mtstate.RequireHookedGeth()

	f.Fuzz(func(t *testing.T, precompileSelector byte, methodSelector byte, input []byte) {
		// Create a StateDB
		sdb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		if err != nil {
			panic(err)
		}
		burner := burn.NewSystemBurner(nil, false)
		_, err = mtosState.InitializeMtosState(sdb, burner, params.MantleDevTestChainConfig())
		if err != nil {
			panic(err)
		}

		// Create an EVM
		gp := core.GasPool(fuzzGas)
		txContext := vm.TxContext{
			GasPrice: common.Big1,
		}
		blockContext := vm.BlockContext{
			CanTransfer: core.CanTransfer,
			Transfer:    core.Transfer,
			GetHash:     nil,
			Coinbase:    common.Address{},
			BlockNumber: new(big.Int),
			Time:        new(big.Int),
			Difficulty:  new(big.Int),
			GasLimit:    fuzzGas,
			BaseFee:     common.Big1,
		}
		evm := vm.NewEVM(blockContext, txContext, sdb, params.MantleDevTestChainConfig(), vm.Config{})

		// Pick a precompile address based on the first byte of the input
		var addr common.Address
		addr[19] = precompileSelector

		// Pick a precompile method based on the second byte of the input
		if precompile := precompiles.Precompiles()[addr]; precompile != nil {
			sigs := precompile.Precompile().Get4ByteMethodSignatures()
			if int(methodSelector) < len(sigs) {
				newInput := make([]byte, 4)
				copy(newInput, sigs[methodSelector][:])
				newInput = append(newInput, input...)
				input = newInput
			}
		}

		// Create and apply a message
		msg := types.NewMessage(
			common.Address{},
			&addr,
			0,
			new(big.Int),
			fuzzGas,
			new(big.Int),
			new(big.Int),
			new(big.Int),
			input,
			nil,
			true,
		)
		_, _ = core.ApplyMessage(evm, msg, &gp)
	})
}
