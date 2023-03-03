package geth

import (
	"bytes"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	evm "github.com/evmos/ethermint/x/evm/vm"
)

var (
	_ evm.EVM         = (*EVM)(nil)
	_ evm.Constructor = NewEVM
)

// EVM is the wrapper for the go-ethereum EVM.
type EVM struct {
	*vm.EVM
}

// NewEVM defines the constructor function for the go-ethereum (geth) EVM. It uses
// the default precompiled contracts and the EVM concrete implementation from
// geth.
func NewEVM(
	ctx sdk.Context,
	blockCtx vm.BlockContext,
	txCtx vm.TxContext,
	stateDB vm.StateDB,
	chainConfig *params.ChainConfig,
	config vm.Config,
	customContracts evm.PrecompiledContracts, // unused
) evm.EVM {
	e := &EVM{
		EVM: vm.NewEVM(blockCtx, txCtx, stateDB, chainConfig, config),
	}

	// pre-compiled contracts
	if len(customContracts) > 0 {
		defaultPrecompiles := vm.DefaultPrecompiles(chainConfig.Rules(blockCtx.BlockNumber, blockCtx.Random != nil))
		active := make([]common.Address, 0, len(customContracts)+len(defaultPrecompiles))
		contracts := make(map[common.Address]vm.PrecompiledContract, len(customContracts)+len(defaultPrecompiles))

		for _, c := range defaultPrecompiles {
			customContracts[c.Address()] = c
			active = append(active, c.Address())
		}

		for _, c := range customContracts {
			if ext, ok := c.(evm.ExtStateDB); ok {
				ext.SetContext(ctx)
				c = ext.(vm.PrecompiledContract)
			}
			contracts[c.Address()] = c
			active = append(active, c.Address())
		}

		sort.SliceStable(active, func(i, j int) bool {
			return bytes.Compare(active[i].Bytes(), active[j].Bytes()) < 0
		})
		e.WithPrecompiles(contracts, active)
	}

	return e
}

// Context returns the EVM's Block Context
func (e EVM) Context() vm.BlockContext {
	return e.EVM.Context
}

// TxContext returns the EVM's Tx Context
func (e EVM) TxContext() vm.TxContext {
	return e.EVM.TxContext
}

// Config returns the configuration options for the EVM.
func (e EVM) Config() vm.Config {
	return e.EVM.Config
}

// Precompile returns the precompiled contract associated with the given address
// and the current chain configuration. If the contract cannot be found it returns
// nil.
func (e EVM) Precompile(addr common.Address) (p vm.PrecompiledContract, found bool) {
	precompiles := GetPrecompiles(e.ChainConfig(), e.EVM.Context.BlockNumber)
	p, found = precompiles[addr]
	return p, found
}

// ActivePrecompiles returns a list of all the active precompiled contract addresses
// for the current chain configuration.
func (EVM) ActivePrecompiles(rules params.Rules) []common.Address {
	return vm.DefaultActivePrecompiles(rules)
}
