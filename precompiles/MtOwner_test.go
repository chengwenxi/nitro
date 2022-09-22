// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/mantlenetworkio/mantle/mtos/burn"
	"github.com/mantlenetworkio/mantle/mtos/mtosState"
	"github.com/mantlenetworkio/mantle/util/testhelpers"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mantlenetworkio/mantle/mtos/util"
)

func TestMtOwner(t *testing.T) {
	evm := newMockEVMForTesting()
	caller := common.BytesToAddress(crypto.Keccak256([]byte{})[:20])
	tracer := util.NewTracingInfo(evm, testhelpers.RandomAddress(), types.MtosAddress, util.TracingDuringEVM)
	state, err := mtosState.OpenMtosState(evm.StateDB, burn.NewSystemBurner(tracer, false))
	Require(t, err)
	Require(t, state.ChainOwners().Add(caller))

	addr1 := common.BytesToAddress(crypto.Keccak256([]byte{1})[:20])
	addr2 := common.BytesToAddress(crypto.Keccak256([]byte{2})[:20])
	addr3 := common.BytesToAddress(crypto.Keccak256([]byte{3})[:20])

	prec := &MtOwner{}
	gasInfo := &MtGasInfo{}
	callCtx := testContext(caller, evm)

	// the zero address is an owner by default
	Require(t, prec.RemoveChainOwner(callCtx, evm, common.Address{}))

	Require(t, prec.AddChainOwner(callCtx, evm, addr1))
	Require(t, prec.AddChainOwner(callCtx, evm, addr2))
	Require(t, prec.AddChainOwner(callCtx, evm, addr1))

	member, err := prec.IsChainOwner(callCtx, evm, addr1)
	Require(t, err)
	if !member {
		Fail(t)
	}

	member, err = prec.IsChainOwner(callCtx, evm, addr2)
	Require(t, err)
	if !member {
		Fail(t)
	}

	member, err = prec.IsChainOwner(callCtx, evm, addr3)
	Require(t, err)
	if member {
		Fail(t)
	}

	Require(t, prec.RemoveChainOwner(callCtx, evm, addr1))
	member, err = prec.IsChainOwner(callCtx, evm, addr1)
	Require(t, err)
	if member {
		Fail(t)
	}
	member, err = prec.IsChainOwner(callCtx, evm, addr2)
	Require(t, err)
	if !member {
		Fail(t)
	}

	Require(t, prec.AddChainOwner(callCtx, evm, addr1))
	all, err := prec.GetAllChainOwners(callCtx, evm)
	Require(t, err)
	if len(all) != 3 {
		Fail(t)
	}
	if all[0] == all[1] || all[1] == all[2] || all[0] == all[2] {
		Fail(t)
	}
	if all[0] != addr1 && all[1] != addr1 && all[2] != addr1 {
		Fail(t)
	}
	if all[0] != addr2 && all[1] != addr2 && all[2] != addr2 {
		Fail(t)
	}
	if all[0] != caller && all[1] != caller && all[2] != caller {
		Fail(t)
	}

	costCap, err := gasInfo.GetAmortizedCostCapBips(callCtx, evm)
	Require(t, err)
	if costCap != math.MaxUint64 {
		Fail(t, costCap)
	}
	newCostCap := uint64(77734)
	Require(t, prec.SetAmortizedCostCapBips(callCtx, evm, newCostCap))
	costCap, err = gasInfo.GetAmortizedCostCapBips(callCtx, evm)
	Require(t, err)
	if costCap != newCostCap {
		Fail(t)
	}
}

func TestMtInfraFeeAccount(t *testing.T) {
	version0 := uint64(0)
	evm := newMockEVMForTestingWithVersion(&version0)
	caller := common.BytesToAddress(crypto.Keccak256([]byte{})[:20])
	newAddr := common.BytesToAddress(crypto.Keccak256([]byte{0})[:20])
	callCtx := testContext(caller, evm)
	prec := &MtOwner{}
	_, err := prec.GetInfraFeeAccount(callCtx, evm)
	Require(t, err)
	err = prec.SetInfraFeeAccount(callCtx, evm, newAddr) // this should be a no-op (because MtOS version 0)
	Require(t, err)

	version5 := uint64(5)
	evm = newMockEVMForTestingWithVersion(&version5)
	callCtx = testContext(caller, evm)
	prec = &MtOwner{}
	precPublic := &MtOwnerPublic{}
	addr, err := prec.GetInfraFeeAccount(callCtx, evm)
	Require(t, err)
	if addr != (common.Address{}) {
		t.Fatal()
	}
	addr, err = precPublic.GetInfraFeeAccount(callCtx, evm)
	Require(t, err)
	if addr != (common.Address{}) {
		t.Fatal()
	}

	err = prec.SetInfraFeeAccount(callCtx, evm, newAddr)
	Require(t, err)
	addr, err = prec.GetInfraFeeAccount(callCtx, evm)
	Require(t, err)
	if addr != newAddr {
		t.Fatal()
	}
	addr, err = precPublic.GetInfraFeeAccount(callCtx, evm)
	Require(t, err)
	if addr != newAddr {
		t.Fatal()
	}
}
