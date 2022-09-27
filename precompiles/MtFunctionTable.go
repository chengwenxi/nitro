// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"errors"
	"math/big"
)

// This precompile provided aggregator's the ability to manage function tables.
// Aggregation works differently in Mantle, so these methods have been stubbed and their effects disabled.
// They are kept for backwards compatibility.
type MtFunctionTable struct {
	Address addr // 0x68
}

// Does nothing
func (con MtFunctionTable) Upload(c ctx, evm mech, buf []byte) error {
	return nil
}

// Returns the empty table's size, which is 0
func (con MtFunctionTable) Size(c ctx, evm mech, addr addr) (huge, error) {
	return big.NewInt(0), nil
}

// Reverts since the table is empty
func (con MtFunctionTable) Get(c ctx, evm mech, addr addr, index huge) (huge, bool, huge, error) {
	return nil, false, nil, errors.New("table is empty")
}
