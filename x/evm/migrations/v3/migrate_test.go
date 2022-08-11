package v3_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	v3 "github.com/evmos/ethermint/x/evm/migrations/v3"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/evmos/ethermint/encoding"

	"github.com/evmos/ethermint/app"
	v3types "github.com/evmos/ethermint/x/evm/migrations/v3/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func TestMigrateStore(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	evmKey := sdk.NewKVStoreKey(types.StoreKey)
	tEvmKey := sdk.NewTransientStoreKey(fmt.Sprintf("%s_test", types.StoreKey))
	ctx := testutil.DefaultContext(evmKey, tEvmKey)
	paramstore := paramtypes.NewSubspace(
		encCfg.Marshaler, encCfg.Amino, evmKey, tEvmKey, "evm",
	).WithKeyTable(v3types.ParamKeyTable())

	params := v3types.DefaultParams()
	paramstore.SetParamSet(ctx, &params)

	require.Panics(t, func() {
		var preMigrationConfig types.ChainConfig
		paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &preMigrationConfig)
	})
	var preMigrationConfig v3types.ChainConfig
	paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &preMigrationConfig)
	require.NotNil(t, preMigrationConfig.MergeForkBlock)

	paramstore = paramtypes.NewSubspace(
		encCfg.Marshaler, encCfg.Amino, evmKey, tEvmKey, "evm",
	).WithKeyTable(types.ParamKeyTable())
	err := v3.MigrateStore(ctx, &paramstore)
	require.NoError(t, err)

	updatedDefaultConfig := types.DefaultChainConfig()

	var postMigrationConfig types.ChainConfig
	paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &postMigrationConfig)
	require.Equal(t, postMigrationConfig.GrayGlacierBlock, updatedDefaultConfig.GrayGlacierBlock)
	require.Equal(t, postMigrationConfig.MergeNetsplitBlock, updatedDefaultConfig.MergeNetsplitBlock)
	require.Panics(t, func() {
		var preMigrationConfig v3types.ChainConfig
		paramstore.Get(ctx, types.ParamStoreKeyChainConfig, &preMigrationConfig)
	})
}

func TestMigrateRejectUnprotectedTx(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	evmKey := sdk.NewKVStoreKey(types.StoreKey)
	tEvmKey := sdk.NewTransientStoreKey(fmt.Sprintf("%s_test", types.StoreKey))
	ctx := testutil.DefaultContext(evmKey, tEvmKey)
	paramsSubspace := paramtypes.NewSubspace(
		encCfg.Marshaler, encCfg.Amino, evmKey, tEvmKey, "evm",
	).WithKeyTable(v3types.ParamKeyTable())

	params := v3types.DefaultParams()
	paramsSubspace.SetParamSet(ctx, &params)

	// check rejectUnprotectedTx
	var rejectUnprotectedTx bool
	require.Panics(t, func() {
		paramsSubspace.Get(ctx, v3types.ParamStoreKeyRejectUnprotectedTx, &rejectUnprotectedTx)
	})

	paramsStore := prefix.NewStore(ctx.KVStore(evmKey), append([]byte(types.ModuleName), '/'))
	// check allowUnprotectedTxs
	bz := paramsStore.Get(types.ParamStoreKeyAllowUnprotectedTxs)
	var allowUnprotectedTxs bool
	err := encCfg.Amino.UnmarshalJSON(bz, &allowUnprotectedTxs)
	require.NoError(t, err)
	require.Equal(t, false, allowUnprotectedTxs)

	// delete allowUnprotectedTxs
	paramsStore.Delete(types.ParamStoreKeyAllowUnprotectedTxs)

	// check allowUnprotectedTxs
	require.Panics(t, func() {
		paramsSubspace.Get(ctx, types.ParamStoreKeyAllowUnprotectedTxs, &allowUnprotectedTxs)
	})

	// set rejectUnprotectedTx
	bz, err = encCfg.Amino.MarshalJSON(false)
	require.NoError(t, err)
	paramsStore.Set(v3types.ParamStoreKeyRejectUnprotectedTx, bz)

	// migrate rejectUnprotectedTx
	err = v3.MigrateRejectUnprotectedTx(ctx, encCfg.Amino, evmKey)
	require.NoError(t, err)

	// check rejectUnprotectedTx
	require.Panics(t, func() {
		paramsSubspace.Get(ctx, v3types.ParamStoreKeyRejectUnprotectedTx, &rejectUnprotectedTx)
	})

	// check allowUnprotectedTxs
	paramsSubspace.Get(ctx, types.ParamStoreKeyAllowUnprotectedTxs, &allowUnprotectedTxs)
	require.Equal(t, true, allowUnprotectedTxs)
}
