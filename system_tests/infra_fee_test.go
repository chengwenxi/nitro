// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantlenetworkio/mantle/blob/main/LICENSE

// race detection makes things slow and miss timeouts
//go:build !race
// +build !race

package mttest

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mantlenetworkio/mantle/mtnode"
	"github.com/mantlenetworkio/mantle/mtos/l2pricing"
	"github.com/mantlenetworkio/mantle/solgen/go/precompilesgen"
	"github.com/mantlenetworkio/mantle/util/mtmath"
)

func TestInfraFee(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodeconfig := mtnode.ConfigDefaultL2Test()

	l2info, _, client, stack := CreateTestL2WithConfig(t, ctx, nil, nodeconfig, true)
	defer requireClose(t, stack)

	l2info.GenerateAccount("User2")

	ownerTxOpts := l2info.GetDefaultTransactOpts("Owner", ctx)
	ownerTxOpts.Context = ctx
	ownerCallOpts := l2info.GetDefaultCallOpts("Owner", ctx)

	mtowner, err := precompilesgen.NewMtOwner(common.HexToAddress("70"), client)
	Require(t, err)
	mtownerPublic, err := precompilesgen.NewMtOwnerPublic(common.HexToAddress("6b"), client)
	Require(t, err)
	networkFeeAddr, err := mtownerPublic.GetNetworkFeeAccount(ownerCallOpts)
	Require(t, err)
	infraFeeAddr := common.BytesToAddress(crypto.Keccak256([]byte{3, 2, 6}))
	tx, err := mtowner.SetInfraFeeAccount(&ownerTxOpts, infraFeeAddr)
	Require(t, err)
	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)

	_, simple := deploySimple(t, ctx, ownerTxOpts, client)

	netFeeBalanceBefore, err := client.BalanceAt(ctx, networkFeeAddr, nil)
	Require(t, err)
	infraFeeBalanceBefore, err := client.BalanceAt(ctx, infraFeeAddr, nil)
	Require(t, err)

	tx, err = simple.Increment(&ownerTxOpts)
	Require(t, err)
	receipt, err := EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)
	l2GasUsed := receipt.GasUsed - receipt.GasUsedForL1
	expectedFunds := mtmath.BigMulByUint(mtmath.UintToBig(l2pricing.InitialBaseFeeWei), l2GasUsed)
	expectedBalanceAfter := mtmath.BigAdd(infraFeeBalanceBefore, expectedFunds)

	netFeeBalanceAfter, err := client.BalanceAt(ctx, networkFeeAddr, nil)
	Require(t, err)
	infraFeeBalanceAfter, err := client.BalanceAt(ctx, infraFeeAddr, nil)
	Require(t, err)

	if !mtmath.BigEquals(netFeeBalanceBefore, netFeeBalanceAfter) {
		Fail(t, netFeeBalanceBefore, netFeeBalanceAfter)
	}
	if !mtmath.BigEquals(infraFeeBalanceAfter, expectedBalanceAfter) {
		Fail(t, infraFeeBalanceBefore, expectedFunds, infraFeeBalanceAfter, expectedBalanceAfter)
	}
}
