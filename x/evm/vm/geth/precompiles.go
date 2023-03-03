package geth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	evm "github.com/evmos/ethermint/x/evm/vm"
)

// GetPrecompiles returns all the precompiled contracts defined given the
// current chain configuration and block height.
func GetPrecompiles(cfg *params.ChainConfig, blockNumber *big.Int) evm.PrecompiledContracts {
	var defaulfPrecompiles evm.PrecompiledContracts
	switch {
	case cfg.IsBerlin(blockNumber):
		defaulfPrecompiles = vm.PrecompiledContractsBerlin
	case cfg.IsIstanbul(blockNumber):
		defaulfPrecompiles = vm.PrecompiledContractsIstanbul
	case cfg.IsByzantium(blockNumber):
		defaulfPrecompiles = vm.PrecompiledContractsByzantium
	default:
		defaulfPrecompiles = vm.PrecompiledContractsHomestead
	}
	precompiles := make(evm.PrecompiledContracts, len(defaulfPrecompiles))
	for address, contract := range defaulfPrecompiles {
		precompiles[address] = contract
	}
	return precompiles
}

// GetActivePrecompiles returns all the precompiled active contracts defined given the
// current chain configuration and block height.
func GetActivePrecompiles(cfg *params.ChainConfig, blockNumber *big.Int) []common.Address {
	var defaultActivePrecompiles []common.Address
	switch {
	case cfg.IsBerlin(blockNumber):
		defaultActivePrecompiles = vm.PrecompiledAddressesBerlin
	case cfg.IsIstanbul(blockNumber):
		defaultActivePrecompiles = vm.PrecompiledAddressesIstanbul
	case cfg.IsByzantium(blockNumber):
		defaultActivePrecompiles = vm.PrecompiledAddressesByzantium
	default:
		defaultActivePrecompiles = vm.PrecompiledAddressesHomestead
	}

	activePrecompiles := make([]common.Address, len(defaultActivePrecompiles))
	copy(activePrecompiles, defaultActivePrecompiles)

	return defaultActivePrecompiles
}
