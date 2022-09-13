// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package precompiles

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/mantlenetworkio/mantle/util/mtmath"
)

// Provides the ability to lookup basic info about accounts and contracts.
type MtInfo struct {
	Address addr // 0x65
}

// Retrieves an account's balance
func (con MtInfo) GetBalance(c ctx, evm mech, account addr) (huge, error) {
	if err := c.Burn(params.BalanceGasEIP1884); err != nil {
		return nil, err
	}
	return evm.StateDB.GetBalance(account), nil
}

// Retrieves a contract's deployed code
func (con MtInfo) GetCode(c ctx, evm mech, account addr) ([]byte, error) {
	if err := c.Burn(params.ColdSloadCostEIP2929); err != nil {
		return nil, err
	}
	code := evm.StateDB.GetCode(account)
	if err := c.Burn(params.CopyGas * mtmath.WordsForBytes(uint64(len(code)))); err != nil {
		return nil, err
	}
	return code, nil
}
