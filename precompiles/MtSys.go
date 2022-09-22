// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/mantlenetworkio/mantle/mtos/util"
	"github.com/mantlenetworkio/mantle/util/merkletree"
	"github.com/mantlenetworkio/mantle/util/mtmath"
)

// Provides system-level functionality for interacting with L1 and understanding the call stack.
type MtSys struct {
	Address                 addr // 0x64
	L2ToL1Tx                func(ctx, mech, addr, addr, huge, huge, huge, huge, huge, huge, []byte) error
	L2ToL1TxGasCost         func(addr, addr, huge, huge, huge, huge, huge, huge, []byte) (uint64, error)
	SendMerkleUpdate        func(ctx, mech, huge, bytes32, huge) error
	SendMerkleUpdateGasCost func(huge, bytes32, huge) (uint64, error)

	// deprecated event
	L2ToL1Transaction        func(ctx, mech, addr, addr, huge, huge, huge, huge, huge, huge, huge, []byte) error
	L2ToL1TransactionGasCost func(addr, addr, huge, huge, huge, huge, huge, huge, huge, []byte) (uint64, error)
}

var InvalidBlockNum = errors.New("Invalid block number")

// Gets the current L2 block number
func (con *MtSys) MtBlockNumber(c ctx, evm mech) (huge, error) {
	return evm.Context.BlockNumber, nil
}

// Gets the L2 block hash, if sufficiently recent
func (con *MtSys) MtBlockHash(c ctx, evm mech, mtBlockNumber *big.Int) (bytes32, error) {
	if !mtBlockNumber.IsUint64() {
		return bytes32{}, InvalidBlockNum
	}
	requestedBlockNum := mtBlockNumber.Uint64()

	currentNumber := evm.Context.BlockNumber.Uint64()
	if requestedBlockNum >= currentNumber || requestedBlockNum+256 < currentNumber {
		return common.Hash{}, errors.New("invalid block number for MtBlockHAsh")
	}

	return evm.Context.GetHash(requestedBlockNum), nil
}

// Gets the rollup's unique chain identifier
func (con *MtSys) MtChainID(c ctx, evm mech) (huge, error) {
	return evm.ChainConfig().ChainID, nil
}

// Gets the current MtOS version
func (con *MtSys) MtOSVersion(c ctx, evm mech) (huge, error) {
	version := new(big.Int).SetUint64(55 + c.State.FormatVersion()) // Mantle starts at version 56
	return version, nil
}

// Returns 0 since Mantle has no concept of storage gas
func (con *MtSys) GetStorageGasAvailable(c ctx, evm mech) (huge, error) {
	return big.NewInt(0), nil
}

// Checks if the call is top-level (deprecated)
func (con *MtSys) IsTopLevelCall(c ctx, evm mech) (bool, error) {
	return evm.Depth() <= 2, nil
}

// Gets the contract's L2 alias
func (con *MtSys) MapL1SenderContractAddressToL2Alias(c ctx, sender addr, dest addr) (addr, error) {
	return util.RemapL1Address(sender), nil
}

// Checks if the caller's caller was aliased
func (con *MtSys) WasMyCallersAddressAliased(c ctx, evm mech) (bool, error) {
	topLevel := con.isTopLevel(c, evm)
	if c.State.FormatVersion() < 6 {
		topLevel = evm.Depth() == 2
	}
	aliased := topLevel && util.DoesTxTypeAlias(c.txProcessor.TopTxType)
	return aliased, nil
}

// Gets the caller's caller without any potential aliasing
func (con *MtSys) MyCallersAddressWithoutAliasing(c ctx, evm mech) (addr, error) {

	address := addr{}

	if evm.Depth() > 1 {
		address = c.txProcessor.Callers[evm.Depth()-2]
	}

	aliased, err := con.WasMyCallersAddressAliased(c, evm)
	if aliased {
		address = util.InverseRemapL1Address(address)
	}
	return address, err
}

// Sends a transaction to L1, adding it to the outbox
func (con *MtSys) SendTxToL1(c ctx, evm mech, value huge, destination addr, calldataForL1 []byte) (huge, error) {
	l1BlockNum, err := c.txProcessor.L1BlockNumber(vm.BlockContext{})
	if err != nil {
		return nil, err
	}
	bigL1BlockNum := mtmath.UintToBig(l1BlockNum)

	mtosState := c.State
	sendHash, err := mtosState.KeccakHash(
		c.caller.Bytes(),
		destination.Bytes(),
		math.U256Bytes(evm.Context.BlockNumber),
		math.U256Bytes(bigL1BlockNum),
		math.U256Bytes(evm.Context.Time),
		common.BigToHash(value).Bytes(),
		calldataForL1,
	)
	if err != nil {
		return nil, err
	}
	merkleAcc := mtosState.SendMerkleAccumulator()
	merkleUpdateEvents, err := merkleAcc.Append(sendHash)
	if err != nil {
		return nil, err
	}

	size, err := merkleAcc.Size()
	if err != nil {
		return nil, err
	}

	// burn the callvalue, which was previously deposited to this precompile's account
	if err := util.BurnBalance(&con.Address, value, evm, util.TracingDuringEVM, "withdraw"); err != nil {
		return nil, err
	}

	for _, merkleUpdateEvent := range merkleUpdateEvents {
		position := merkletree.LevelAndLeaf{
			Level: merkleUpdateEvent.Level,
			Leaf:  merkleUpdateEvent.NumLeaves,
		}
		err := con.SendMerkleUpdate(
			c,
			evm,
			big.NewInt(0),
			merkleUpdateEvent.Hash,
			position.ToBigInt(),
		)
		if err != nil {
			return nil, err
		}
	}

	leafNum := big.NewInt(int64(size - 1))

	err = con.L2ToL1Tx(
		c,
		evm,
		c.caller,
		destination,
		sendHash.Big(),
		leafNum,
		evm.Context.BlockNumber,
		bigL1BlockNum,
		evm.Context.Time,
		value,
		calldataForL1,
	)

	if c.State.FormatVersion() >= 4 {
		return leafNum, nil
	}
	return sendHash.Big(), err
}

// Gets the root, size, and partials of the outbox Merkle tree state (caller must be the 0 address)
func (con MtSys) SendMerkleTreeState(c ctx, evm mech) (huge, bytes32, []bytes32, error) {
	if c.caller != (addr{}) {
		return nil, bytes32{}, nil, errors.New("method can only be called by address zero")
	}

	// OK to not charge gas, because method is only callable by address zero

	size, rootHash, rawPartials, _ := c.State.SendMerkleAccumulator().StateForExport()
	partials := make([]bytes32, len(rawPartials))
	for i, par := range rawPartials {
		partials[i] = bytes32(par)
	}
	return big.NewInt(int64(size)), bytes32(rootHash), partials, nil
}

// Send paid eth to the destination on L1
func (con MtSys) WithdrawEth(c ctx, evm mech, value huge, destination addr) (huge, error) {
	return con.SendTxToL1(c, evm, value, destination, []byte{})
}

func (con MtSys) isTopLevel(c ctx, evm mech) bool {
	depth := evm.Depth()
	return depth < 2 || evm.Origin == c.txProcessor.Callers[depth-2]
}
