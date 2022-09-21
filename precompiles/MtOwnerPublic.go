// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package precompiles

import (
	"github.com/ethereum/go-ethereum/common"
)

// This precompile provides non-owners with info about the current chain owners.
// The calls to this precompile do not require the sender be a chain owner.
// For those that are, see MtOwner
type MtOwnerPublic struct {
	Address addr // 0x6b
}

// Retrieves the list of chain owners
func (con MtOwnerPublic) GetAllChainOwners(c ctx, evm mech) ([]common.Address, error) {
	return c.State.ChainOwners().AllMembers(65536)
}

// See if the user is a chain owner
func (con MtOwnerPublic) IsChainOwner(c ctx, evm mech, addr addr) (bool, error) {
	return c.State.ChainOwners().IsMember(addr)
}

// Gets the network fee collector
func (con MtOwnerPublic) GetNetworkFeeAccount(c ctx, evm mech) (addr, error) {
	return c.State.NetworkFeeAccount()
}

// Gets the infrastructure fee collector
func (con MtOwnerPublic) GetInfraFeeAccount(c ctx, evm mech) (addr, error) {
	if c.State.FormatVersion() < 6 {
		return c.State.NetworkFeeAccount()
	}
	return c.State.InfraFeeAccount()
}
