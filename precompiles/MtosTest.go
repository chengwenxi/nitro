// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"errors"
)

// Provides a method of burning arbitrary amounts of gas, which exists for historical reasons.
type MtosTest struct {
	Address addr // 0x69
}

// Unproductively burns the amount of L2 MtGas
func (con MtosTest) BurnMtGas(c ctx, gasAmount huge) error {
	if !gasAmount.IsUint64() {
		return errors.New("Not a uint64")
	}
	//nolint:errcheck
	c.Burn(gasAmount.Uint64()) // burn the amount, even if it's more than the user has
	return nil
}
