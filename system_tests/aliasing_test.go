// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mttest

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mantlenetworkio/mantle/mtos/util"
	"github.com/mantlenetworkio/mantle/solgen/go/mocksgen"
	"github.com/mantlenetworkio/mantle/solgen/go/precompilesgen"
)

func TestAliasing(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l2info, _, l2client, l2stack, l1info, _, l1client, l1stack := createTestNodeOnL1(t, ctx, true)
	defer requireClose(t, l1stack)
	defer requireClose(t, l2stack)

	auth := l2info.GetDefaultTransactOpts("Owner", ctx)
	user := l1info.GetDefaultTransactOpts("User", ctx)
	TransferBalanceTo(t, "Owner", util.RemapL1Address(user.From), big.NewInt(1e18), l2info, l2client, ctx)

	simpleAddr, simple := deploySimple(t, ctx, auth, l2client)
	simpleContract, err := abi.JSON(strings.NewReader(mocksgen.SimpleABI))
	Require(t, err)

	// Test direct calls
	mtsys, err := precompilesgen.NewMtSys(types.MtSysAddress, l2client)
	Require(t, err)
	top, err := mtsys.IsTopLevelCall(nil)
	Require(t, err)
	was, err := mtsys.WasMyCallersAddressAliased(nil)
	Require(t, err)
	alias, err := mtsys.MyCallersAddressWithoutAliasing(nil)
	Require(t, err)
	if !top {
		Fail(t, "direct call is not top level")
	}
	if was || alias != (common.Address{}) {
		Fail(t, "direct call has an alias", was, alias)
	}

	testL2Signed := func(top, direct, static, delegate, callcode, call bool) {
		t.Helper()

		// check via L2
		tx, err := simple.CheckCalls(&auth, top, direct, static, delegate, callcode, call)
		Require(t, err)
		_, err = EnsureTxSucceeded(ctx, l2client, tx)
		Require(t, err)

		// check signed txes via L1
		data, err := simpleContract.Pack("checkCalls", top, direct, static, delegate, callcode, call)
		Require(t, err)
		tx = l2info.PrepareTxTo("Owner", &simpleAddr, 500000, big.NewInt(0), data)
		SendSignedTxViaL1(t, ctx, l1info, l1client, l2client, tx)
	}

	testUnsigned := func(top, direct, static, delegate, callcode, call bool) {
		t.Helper()

		// check unsigned txes via L1
		data, err := simpleContract.Pack("checkCalls", top, direct, static, delegate, callcode, call)
		Require(t, err)
		tx := l2info.PrepareTxTo("Owner", &simpleAddr, 500000, big.NewInt(0), data)
		SendUnsignedTxViaL1(t, ctx, l1info, l1client, l2client, tx)
	}

	testL2Signed(true, true, false, false, false, false)
	testL2Signed(false, false, false, false, false, false)
	testUnsigned(true, true, false, false, false, false)
	testUnsigned(false, true, false, true, false, false)
}
