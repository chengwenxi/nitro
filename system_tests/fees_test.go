// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

// these tests seems to consume too much memory with race detection
//go:build !race
// +build !race

package mttest

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mantlenetworkio/mantle/mtcompress"
	"github.com/mantlenetworkio/mantle/mtnode"
	"github.com/mantlenetworkio/mantle/mtos/l1pricing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mantlenetworkio/mantle/solgen/go/precompilesgen"
	"github.com/mantlenetworkio/mantle/util/colors"
	"github.com/mantlenetworkio/mantle/util/mtmath"
)

func TestSequencerFeePaid(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l2info, _, l2client, l2stack, _, _, _, l1stack := createTestNodeOnL1(t, ctx, true)
	defer requireClose(t, l1stack)
	defer requireClose(t, l2stack)

	callOpts := l2info.GetDefaultCallOpts("Owner", ctx)

	// get the network fee account
	arbOwnerPublic, err := precompilesgen.NewMtOwnerPublic(common.HexToAddress("0x6b"), l2client)
	Require(t, err, "could not deploy MtOwner contract")
	arbGasInfo, err := precompilesgen.NewMtGasInfo(common.HexToAddress("0x6c"), l2client)
	Require(t, err, "could not deploy MtOwner contract")
	networkFeeAccount, err := arbOwnerPublic.GetNetworkFeeAccount(callOpts)
	Require(t, err, "could not get the network fee account")

	l1Estimate, err := arbGasInfo.GetL1BaseFeeEstimate(callOpts)
	Require(t, err)
	networkBefore := GetBalance(t, ctx, l2client, networkFeeAccount)

	l2info.GasPrice = GetBaseFee(t, l2client, ctx)
	tx, receipt := TransferBalance(t, "Faucet", "Faucet", big.NewInt(0), l2info, l2client, ctx)
	txSize := compressedTxSize(t, tx)

	networkAfter := GetBalance(t, ctx, l2client, networkFeeAccount)
	l1Charge := mtmath.BigMulByUint(l2info.GasPrice, receipt.GasUsedForL1)

	networkRevenue := mtmath.BigSub(networkAfter, networkBefore)
	gasUsedForL2 := receipt.GasUsed - receipt.GasUsedForL1
	if !mtmath.BigEquals(networkRevenue, mtmath.BigMulByUint(tx.GasPrice(), gasUsedForL2)) {
		Fail(t, "network didn't receive expected payment")
	}

	l1GasBought := mtmath.BigDiv(l1Charge, l1Estimate).Uint64()
	l1GasActual := txSize * params.TxDataNonZeroGasEIP2028

	colors.PrintBlue("bytes ", l1GasBought/params.TxDataNonZeroGasEIP2028, txSize)

	if l1GasBought != l1GasActual {
		Fail(t, "the sequencer's future revenue does not match its costs", l1GasBought, l1GasActual)
	}
}

func testSequencerPriceAdjustsFrom(t *testing.T, initialEstimate uint64) {
	t.Parallel()

	_ = os.Mkdir("test-data", 0766)
	path := filepath.Join("test-data", fmt.Sprintf("testSequencerPriceAdjustsFrom%v.csv", initialEstimate))

	f, err := os.Create(path)
	Require(t, err)
	defer func() { Require(t, f.Close()) }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chainConfig := params.MantleDevTestChainConfig()
	conf := mtnode.ConfigDefaultL1Test()
	conf.DelayedSequencer.FinalizeDistance = 1

	l2info, node, l2client, l2stack, _, _, l1client, l1stack := createTestNodeOnL1WithConfig(t, ctx, true, conf, chainConfig)
	defer requireClose(t, l1stack)
	defer requireClose(t, l2stack)

	ownerAuth := l2info.GetDefaultTransactOpts("Owner", ctx)

	// make ownerAuth a chain owner
	arbdebug, err := precompilesgen.NewMtDebug(common.HexToAddress("0xff"), l2client)
	Require(t, err)
	tx, err := arbdebug.BecomeChainOwner(&ownerAuth)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, l2client, tx)

	// use ownerAuth to set the L1 price per unit
	Require(t, err)
	arbOwner, err := precompilesgen.NewMtOwner(common.HexToAddress("0x70"), l2client)
	Require(t, err)
	tx, err = arbOwner.SetL1PricePerUnit(&ownerAuth, mtmath.UintToBig(initialEstimate))
	Require(t, err)
	_, err = WaitForTx(ctx, l2client, tx.Hash(), time.Second*5)
	Require(t, err)

	arbGasInfo, err := precompilesgen.NewMtGasInfo(common.HexToAddress("0x6c"), l2client)
	Require(t, err)
	lastEstimate, err := arbGasInfo.GetL1BaseFeeEstimate(&bind.CallOpts{Context: ctx})
	Require(t, err)
	lastBatchCount, err := node.InboxTracker.GetBatchCount()
	Require(t, err)
	l1Header, err := l1client.HeaderByNumber(ctx, nil)
	Require(t, err)

	rewardRecipientBalanceBefore := GetBalance(t, ctx, l2client, l1pricing.BatchPosterAddress)
	timesPriceAdjusted := 0

	colors.PrintBlue("Initial values")
	colors.PrintBlue("    L1 base fee ", l1Header.BaseFee)
	colors.PrintBlue("    L1 estimate ", lastEstimate)

	numRetrogradeMoves := 0
	for i := 0; i < 256; i++ {
		tx, receipt := TransferBalance(t, "Owner", "Owner", common.Big1, l2info, l2client, ctx)
		header, err := l2client.HeaderByHash(ctx, receipt.BlockHash)
		Require(t, err)

		units := compressedTxSize(t, tx) * params.TxDataNonZeroGasEIP2028
		estimatedL1FeePerUnit := mtmath.BigDivByUint(mtmath.BigMulByUint(header.BaseFee, receipt.GasUsedForL1), units)

		if !mtmath.BigEquals(lastEstimate, estimatedL1FeePerUnit) {
			l1Header, err = l1client.HeaderByNumber(ctx, nil)
			Require(t, err)

			callOpts := &bind.CallOpts{Context: ctx, BlockNumber: receipt.BlockNumber}
			actualL1FeePerUnit, err := arbGasInfo.GetL1BaseFeeEstimate(callOpts)
			Require(t, err)
			surplus, err := arbGasInfo.GetL1PricingSurplus(callOpts)
			Require(t, err)

			colors.PrintGrey("MtOS updated its L1 estimate")
			colors.PrintGrey("    L1 base fee ", l1Header.BaseFee)
			colors.PrintGrey("    L1 estimate ", lastEstimate, " ➤ ", estimatedL1FeePerUnit, " = ", actualL1FeePerUnit)
			colors.PrintGrey("    Surplus ", surplus)
			fmt.Fprintf(
				f, "%v, %v, %v, %v, %v, %v\n", i, l1Header.BaseFee, lastEstimate,
				estimatedL1FeePerUnit, actualL1FeePerUnit, surplus,
			)

			oldDiff := mtmath.BigAbs(mtmath.BigSub(lastEstimate, l1Header.BaseFee))
			newDiff := mtmath.BigAbs(mtmath.BigSub(actualL1FeePerUnit, l1Header.BaseFee))
			cmpDiff := mtmath.BigGreaterThan(newDiff, oldDiff)
			signums := surplus.Sign() == mtmath.BigSub(actualL1FeePerUnit, l1Header.BaseFee).Sign()

			if timesPriceAdjusted > 0 && cmpDiff && signums {
				numRetrogradeMoves++
				if numRetrogradeMoves > 1 {
					colors.PrintRed(timesPriceAdjusted, newDiff, oldDiff, lastEstimate, surplus)
					colors.PrintRed(estimatedL1FeePerUnit, l1Header.BaseFee, actualL1FeePerUnit)
					Fail(t, "L1 gas price estimate should tend toward the basefee")
				}
			} else {
				numRetrogradeMoves = 0
			}
			diff := mtmath.BigAbs(mtmath.BigSub(actualL1FeePerUnit, estimatedL1FeePerUnit))
			maxDiffToAllow := mtmath.BigDivByUint(actualL1FeePerUnit, 100)
			if mtmath.BigLessThan(maxDiffToAllow, diff) { // verify that estimates is within 1% of actual
				Fail(t, "New L1 estimate differs too much from receipt")
			}
			if mtmath.BigEquals(actualL1FeePerUnit, common.Big0) {
				Fail(t, "Estimate is zero", i)
			}
			lastEstimate = actualL1FeePerUnit
			timesPriceAdjusted++
		}

		if i%16 == 0 {
			// see that the inbox advances

			for j := 16; j > 0; j-- {
				newBatchCount, err := node.InboxTracker.GetBatchCount()
				Require(t, err)
				if newBatchCount > lastBatchCount {
					colors.PrintGrey("posted new batch ", newBatchCount)
					lastBatchCount = newBatchCount
					break
				}
				if j == 1 {
					Fail(t, "batch count didn't update in time")
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}

	rewardRecipientBalanceAfter := GetBalance(t, ctx, l2client, chainConfig.MantleChainParams.InitialChainOwner)
	colors.PrintMint("reward recipient balance ", rewardRecipientBalanceBefore, " ➤ ", rewardRecipientBalanceAfter)
	colors.PrintMint("price changes     ", timesPriceAdjusted)

	if timesPriceAdjusted == 0 {
		Fail(t, "L1 gas price estimate never adjusted")
	}
	if !mtmath.BigGreaterThan(rewardRecipientBalanceAfter, rewardRecipientBalanceBefore) {
		Fail(t, "reward recipient didn't get paid")
	}

	arbAggregator, err := precompilesgen.NewMtAggregator(common.HexToAddress("0x6d"), l2client)
	Require(t, err)
	batchPosterAddresses, err := arbAggregator.GetBatchPosters(&bind.CallOpts{Context: ctx})
	Require(t, err)
	numReimbursed := 0
	for _, bpAddr := range batchPosterAddresses {
		if bpAddr != l1pricing.BatchPosterAddress && bpAddr != l1pricing.L1PricerFundsPoolAddress {
			numReimbursed++
			bal, err := l1client.BalanceAt(ctx, bpAddr, nil)
			Require(t, err)
			if bal.Sign() == 0 {
				Fail(t, "Batch poster balance is zero for", bpAddr)
			}
		}
	}
	if numReimbursed != 1 {
		Fail(t, "Wrong number of batch posters were reimbursed", numReimbursed)
	}
}

func TestSequencerPriceAdjustsFrom1Gwei(t *testing.T) {
	testSequencerPriceAdjustsFrom(t, params.GWei)
}

func TestSequencerPriceAdjustsFrom2Gwei(t *testing.T) {
	testSequencerPriceAdjustsFrom(t, 2*params.GWei)
}

func TestSequencerPriceAdjustsFrom5Gwei(t *testing.T) {
	testSequencerPriceAdjustsFrom(t, 5*params.GWei)
}

func TestSequencerPriceAdjustsFrom10Gwei(t *testing.T) {
	testSequencerPriceAdjustsFrom(t, 10*params.GWei)
}

func TestSequencerPriceAdjustsFrom25Gwei(t *testing.T) {
	testSequencerPriceAdjustsFrom(t, 25*params.GWei)
}

func compressedTxSize(t *testing.T, tx *types.Transaction) uint64 {
	txBin, err := tx.MarshalBinary()
	Require(t, err)
	compressed, err := mtcompress.CompressFast(txBin)
	Require(t, err)
	return uint64(len(compressed))
}
