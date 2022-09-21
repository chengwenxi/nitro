// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"errors"
	"math/big"
)

// This precompile provides the ability to create short-hands for commonly used accounts.
type MtAddressTable struct {
	Address addr // 0x66
}

// Check if an address exists in the table
func (con MtAddressTable) AddressExists(c ctx, evm mech, addr addr) (bool, error) {
	return c.State.AddressTable().AddressExists(addr)
}

// Gets bytes that represent the address
func (con MtAddressTable) Compress(c ctx, evm mech, addr addr) ([]uint8, error) {
	return c.State.AddressTable().Compress(addr)
}

// Replaces the compressed bytes at the given offset with those of the corresponding account
func (con MtAddressTable) Decompress(c ctx, evm mech, buf []uint8, offset huge) (addr, huge, error) {
	if !offset.IsInt64() {
		return addr{}, nil, errors.New("invalid offset in MtAddressTable.Decompress")
	}
	ioffset := offset.Int64()
	if ioffset > int64(len(buf)) {
		return addr{}, nil, errors.New("invalid offset in MtAddressTable.Decompress")
	}
	result, nbytes, err := c.State.AddressTable().Decompress(buf[ioffset:])
	return result, big.NewInt(int64(nbytes)), err
}

// Looks up the index of an address in the table
func (con MtAddressTable) Lookup(c ctx, evm mech, addr addr) (huge, error) {
	result, exists, err := c.State.AddressTable().Lookup(addr)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("address does not exist in AddressTable")
	}
	return big.NewInt(int64(result)), nil
}

// Looks up an address in the table by index
func (con MtAddressTable) LookupIndex(c ctx, evm mech, index huge) (addr, error) {
	if !index.IsUint64() {
		return addr{}, errors.New("invalid index in MtAddressTable.LookupIndex")
	}
	result, exists, err := c.State.AddressTable().LookupIndex(index.Uint64())
	if err != nil {
		return addr{}, err
	}
	if !exists {
		return addr{}, errors.New("index does not exist in AddressTable")
	}
	return result, nil
}

// Adds an account to the table, shrinking its compressed representation
func (con MtAddressTable) Register(c ctx, evm mech, addr addr) (huge, error) {
	slot, err := c.State.AddressTable().Register(addr)
	return big.NewInt(int64(slot)), err
}

// Gets the number of addresses in the table
func (con MtAddressTable) Size(c ctx, evm mech) (huge, error) {
	size, err := c.State.AddressTable().Size()
	return big.NewInt(int64(size)), err
}
