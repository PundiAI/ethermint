package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v3types "github.com/evmos/ethermint/x/evm/migrations/v3/types"
	"github.com/evmos/ethermint/x/evm/types"
)

// MigrateStore sets the default for GrayGlacierBlock and MergeNetsplitBlock in ChainConfig parameter.
func MigrateStore(ctx sdk.Context, paramstore *paramtypes.Subspace) error {
	if !paramstore.HasKeyTable() {
		ps := paramstore.WithKeyTable(types.ParamKeyTable())
		paramstore = &ps
	}
	prevConfig := &types.ChainConfig{}
	paramstore.GetIfExists(ctx, types.ParamStoreKeyChainConfig, prevConfig)

	defaultConfig := types.DefaultChainConfig()

	prevConfig.GrayGlacierBlock = defaultConfig.GrayGlacierBlock
	prevConfig.MergeNetsplitBlock = defaultConfig.MergeNetsplitBlock

	paramstore.Set(ctx, types.ParamStoreKeyChainConfig, prevConfig)
	return nil
}

// MigrateRejectUnprotectedTx used by ethermint version before v0.17.0
func MigrateRejectUnprotectedTx(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) error {
	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(types.ModuleName), '/'))

	bzR := paramsStore.Get(v3types.ParamStoreKeyRejectUnprotectedTx)
	var rejectUnprotectedTx bool
	if err := legacyAmino.UnmarshalJSON(bzR, &rejectUnprotectedTx); err != nil {
		return fmt.Errorf("legacy amino unmarshal %s: %s", err.Error(), v3types.ParamStoreKeyRejectUnprotectedTx)
	}

	allowUnprotectedTxs := !rejectUnprotectedTx
	bzA, err := legacyAmino.MarshalJSON(allowUnprotectedTxs)
	if err != nil {
		return fmt.Errorf("legacy amino marshal %s: %s", err.Error(), v3types.ParamStoreKeyRejectUnprotectedTx)
	}

	ctx.Logger().Info("migrate params", "module", types.ModuleName, "from", fmt.Sprintf("%s:%v", v3types.ParamStoreKeyRejectUnprotectedTx, rejectUnprotectedTx),
		"to", fmt.Sprintf("%s:%v", types.ParamStoreKeyAllowUnprotectedTxs, allowUnprotectedTxs))

	paramsStore.Delete(v3types.ParamStoreKeyRejectUnprotectedTx)
	paramsStore.Set(types.ParamStoreKeyAllowUnprotectedTxs, bzA)
	return nil
}
