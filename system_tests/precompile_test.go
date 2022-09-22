// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mttest

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mantlenetworkio/mantle/solgen/go/precompilesgen"
)

func TestPurePrecompileMethodCalls(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, client, l2stack := CreateTestL2(t, ctx)
	defer requireClose(t, l2stack)

	mtSys, err := precompilesgen.NewMtSys(common.HexToAddress("0x64"), client)
	Require(t, err, "could not deploy MtSys contract")
	chainId, err := mtSys.MtChainID(&bind.CallOpts{})
	Require(t, err, "failed to get the ChainID")
	if chainId.Uint64() != params.MantleDevTestChainConfig().ChainID.Uint64() {
		Fail(t, "Wrong ChainID", chainId.Uint64())
	}
}

func TestCustomSolidityErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, client, l2stack := CreateTestL2(t, ctx)
	defer requireClose(t, l2stack)

	mtDebug, err := precompilesgen.NewMtDebug(common.HexToAddress("0xff"), client)
	Require(t, err, "could not deploy MtDebug contract")
	customError := mtDebug.CustomRevert(&bind.CallOpts{}, 1024)
	if customError == nil {
		Fail(t, "should have errored")
	}
	observedMessage := customError.Error()
	expectedMessage := "execution reverted: error Custom(1024, This spider family wards off bugs: /\\oo/\\ //\\(oo)/\\ /\\oo/\\, true)"
	if observedMessage != expectedMessage {
		Fail(t, observedMessage)
	}
}
