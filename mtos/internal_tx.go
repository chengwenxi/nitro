// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package mtos

import (
	"fmt"
	"math/big"

	"github.com/mantlenetworkio/mantle/util/mtmath"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/mantlenetworkio/mantle/mtos/mtosState"
	"github.com/mantlenetworkio/mantle/mtos/util"
)

func InternalTxStartBlock(
	chainId,
	l1BaseFee *big.Int,
	l1BlockNum uint64,
	header,
	lastHeader *types.Header,
) *types.MantleInternalTx {

	l2BlockNum := header.Number.Uint64()
	timePassed := header.Time - lastHeader.Time

	if l1BaseFee == nil {
		l1BaseFee = big.NewInt(0)
	}
	data, err := util.PackInternalTxDataStartBlock(l1BaseFee, l1BlockNum, l2BlockNum, timePassed)
	if err != nil {
		panic(fmt.Sprintf("Failed to pack internal tx %v", err))
	}
	return &types.MantleInternalTx{
		ChainId: chainId,
		Data:    data,
	}
}

func ApplyInternalTxUpdate(tx *types.MantleInternalTx, state *mtosState.MtosState, evm *vm.EVM) {
	switch *(*[4]byte)(tx.Data[:4]) {
	case InternalTxStartBlockMethodID:
		inputs, err := util.UnpackInternalTxDataStartBlock(tx.Data)
		if err != nil {
			panic(err)
		}

		l1BlockNumber := util.SafeMapGet[uint64](inputs, "l1BlockNumber")
		timePassed := util.SafeMapGet[uint64](inputs, "timePassed")
		if state.FormatVersion() < 3 {
			// (incorrectly) use the L2 block number instead
			timePassed = util.SafeMapGet[uint64](inputs, "l2BlockNumber")
		}

		nextL1BlockNumber, err := state.Blockhashes().NextBlockNumber()
		state.Restrict(err)

		l2BaseFee, err := state.L2PricingState().BaseFeeWei()
		state.Restrict(err)

		if l1BlockNumber >= nextL1BlockNumber {
			var prevHash common.Hash
			if evm.Context.BlockNumber.Sign() > 0 {
				prevHash = evm.Context.GetHash(evm.Context.BlockNumber.Uint64() - 1)
			}
			state.Restrict(state.Blockhashes().RecordNewL1Block(l1BlockNumber, prevHash))
		}

		currentTime := evm.Context.Time.Uint64()

		// Try to reap 2 retryables
		_ = state.RetryableState().TryToReapOneRetryable(currentTime, evm, util.TracingDuringEVM)
		_ = state.RetryableState().TryToReapOneRetryable(currentTime, evm, util.TracingDuringEVM)

		state.L2PricingState().UpdatePricingModel(l2BaseFee, timePassed, false)

		state.UpgradeMtosVersionIfNecessary(currentTime)
	case InternalTxBatchPostingReportMethodID:
		inputs, err := util.UnpackInternalTxDataBatchPostingReport(tx.Data)
		if err != nil {
			panic(err)
		}
		batchTimestamp := util.SafeMapGet[*big.Int](inputs, "batchTimestamp")
		batchPosterAddress := util.SafeMapGet[common.Address](inputs, "batchPosterAddress")
		batchDataGas := util.SafeMapGet[uint64](inputs, "batchDataGas")
		l1BaseFeeWei := util.SafeMapGet[*big.Int](inputs, "l1BaseFeeWei")

		l1p := state.L1PricingState()
		perBatchGas, err := l1p.PerBatchGasCost()
		if err != nil {
			log.Warn("L1Pricing PerBatchGas failed", "err", err)
		}
		gasSpent := mtmath.SaturatingAdd(perBatchGas, mtmath.SaturatingCast(batchDataGas))
		weiSpent := mtmath.BigMulByUint(l1BaseFeeWei, mtmath.SaturatingUCast(gasSpent))
		err = l1p.UpdateForBatchPosterSpending(
			evm.StateDB,
			evm,
			state.FormatVersion(),
			batchTimestamp.Uint64(),
			evm.Context.Time.Uint64(),
			batchPosterAddress,
			weiSpent,
			l1BaseFeeWei,
			util.TracingDuringEVM,
		)
		if err != nil {
			log.Warn("L1Pricing UpdateForSequencerSpending failed", "err", err)
		}
	}
}
