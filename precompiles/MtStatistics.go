// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"math/big"
)

// Provides statistics about the rollup right before the Mantle upgrade.
// In Classic, this was how a user would get info such as the total number of accounts,
// but there's now better ways to do that with geth.
type MtStatistics struct {
	Address addr // 0x6e
}

// Returns the current block number and some statistics about the rollup's pre-Mantle state
func (con MtStatistics) GetStats(c ctx, evm mech) (huge, huge, huge, huge, huge, huge, error) {
	blockNum := evm.Context.BlockNumber
	classicNumAccounts := big.NewInt(0)  // TODO: hardcode the final value from Mantle Classic
	classicStorageSum := big.NewInt(0)   // TODO: hardcode the final value from Mantle Classic
	classicGasSum := big.NewInt(0)       // TODO: hardcode the final value from Mantle Classic
	classicNumTxes := big.NewInt(0)      // TODO: hardcode the final value from Mantle Classic
	classicNumContracts := big.NewInt(0) // TODO: hardcode the final value from Mantle Classic
	return blockNum, classicNumAccounts, classicStorageSum, classicGasSum, classicNumTxes, classicNumContracts, nil
}
