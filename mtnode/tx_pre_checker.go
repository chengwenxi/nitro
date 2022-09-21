// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mtnode

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mantlenetworkio/mantle/mtos/l1pricing"
	"github.com/mantlenetworkio/mantle/mtos/mtosState"
	"github.com/mantlenetworkio/mantle/util/mtmath"
)

type TxPreChecker struct {
	TransactionPublisher
	bc            *core.BlockChain
	getStrictness func() uint
}

func NewTxPreChecker(publisher TransactionPublisher, bc *core.BlockChain, getStrictness func() uint) *TxPreChecker {
	return &TxPreChecker{
		TransactionPublisher: publisher,
		bc:                   bc,
		getStrictness:        getStrictness,
	}
}

const TxPreCheckerStrictnessNone uint = 0
const TxPreCheckerStrictnessAlwaysCompatible uint = 10
const TxPreCheckerStrictnessLikelyCompatible uint = 20
const TxPreCheckerStrictnessFullValidation uint = 30

func PreCheckTx(chainConfig *params.ChainConfig, header *types.Header, statedb *state.StateDB, mtos *mtosState.MtosState, tx *types.Transaction, strictness uint) error {
	if strictness < TxPreCheckerStrictnessAlwaysCompatible {
		return nil
	}
	if tx.Gas() < params.TxGas {
		return core.ErrIntrinsicGas
	}
	sender, err := types.Sender(types.MakeSigner(chainConfig, header.Number), tx)
	if err != nil {
		return err
	}
	baseFee := header.BaseFee
	if strictness < TxPreCheckerStrictnessLikelyCompatible {
		baseFee, err = mtos.L2PricingState().MinBaseFeeWei()
		if err != nil {
			return err
		}
	}
	if mtmath.BigLessThan(tx.GasFeeCap(), baseFee) {
		return fmt.Errorf("%w: address %v, maxFeePerGas: %s baseFee: %s", core.ErrFeeCapTooLow, sender, tx.GasFeeCap(), header.BaseFee)
	}
	stateNonce := statedb.GetNonce(sender)
	if tx.Nonce() < stateNonce {
		return fmt.Errorf("%w: address %v, tx: %d state: %d", core.ErrNonceTooLow, sender, tx.Nonce(), stateNonce)
	}
	intrinsic, err := core.IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, chainConfig.IsHomestead(header.Number), true)
	if err != nil {
		return err
	}
	if tx.Gas() < intrinsic {
		return core.ErrIntrinsicGas
	}
	if strictness < TxPreCheckerStrictnessLikelyCompatible {
		return nil
	}
	balance := statedb.GetBalance(sender)
	cost := tx.Cost()
	if mtmath.BigLessThan(balance, cost) {
		return fmt.Errorf("%w: address %v have %v want %v", core.ErrInsufficientFunds, sender, balance, cost)
	}
	if strictness >= TxPreCheckerStrictnessFullValidation && tx.Nonce() > stateNonce {
		return fmt.Errorf("%w: address %v, tx: %d state: %d", core.ErrNonceTooHigh, sender, tx.Nonce(), stateNonce)
	}
	dataCost, _ := mtos.L1PricingState().GetPosterInfo(tx, l1pricing.BatchPosterAddress)
	dataGas := mtmath.BigDiv(dataCost, header.BaseFee)
	if tx.Gas() < intrinsic+dataGas.Uint64() {
		return core.ErrIntrinsicGas
	}
	return nil
}

func (c *TxPreChecker) PublishTransaction(ctx context.Context, tx *types.Transaction) error {
	block := c.bc.CurrentBlock()
	statedb, err := c.bc.StateAt(block.Root())
	if err != nil {
		return err
	}
	mtos, err := mtosState.OpenSystemMtosState(statedb, nil, true)
	if err != nil {
		return err
	}
	err = PreCheckTx(c.bc.Config(), block.Header(), statedb, mtos, tx, c.getStrictness())
	if err != nil {
		return err
	}
	return c.TransactionPublisher.PublishTransaction(ctx, tx)
}
