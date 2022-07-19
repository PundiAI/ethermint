package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	// TODO: update once https://github.com/cosmos/cosmos-sdk/pull/12615 is merged
	// and released
	for _, pair := range params.ParamSetPairs() {
		k.paramSpace.GetIfExists(ctx, pair.Key, pair.Value)
	}

	k.paramSpace.GetParamSet(ctx, &params)
	// TODO params store RejectUnprotectedTx, value is the opposite of AllowUnprotectedTxs
	params.AllowUnprotectedTxs = !params.AllowUnprotectedTxs
	return params
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	// NOTE: params AllowUnprotectedTxs
	k.paramSpace.SetParamSet(ctx, &params)
}
