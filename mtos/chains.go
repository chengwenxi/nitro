// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mtos

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func getStaticChainConfig(chainId *big.Int) (*params.ChainConfig, error) {
	for _, potentialChainConfig := range params.MantleSupportedChainConfigs {
		if potentialChainConfig.ChainID.Cmp(chainId) == 0 {
			return potentialChainConfig, nil
		}
	}
	return nil, fmt.Errorf("unsupported L2 chain ID %v", chainId)
}

func GetChainConfig(chainId *big.Int, genesisBlockNum uint64) (*params.ChainConfig, error) {
	staticChainConfig, err := getStaticChainConfig(chainId)
	if err != nil {
		return nil, err
	}
	staticChainConfig.MantleChainParams.GenesisBlockNum = genesisBlockNum
	return staticChainConfig, nil
}
