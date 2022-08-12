package keeper_test

import (
	"github.com/evmos/ethermint/x/evm/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.EvmKeeper.GetParams(suite.ctx)
	defaultParams := types.DefaultParams()
	defaultParams.AllowUnprotectedTxs = !defaultParams.AllowUnprotectedTxs
	suite.Require().Equal(defaultParams, params)
	params.EvmDenom = "inj"
	params.AllowUnprotectedTxs = !params.AllowUnprotectedTxs // NOTE
	suite.app.EvmKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.EvmKeeper.GetParams(suite.ctx)
	params.AllowUnprotectedTxs = !params.AllowUnprotectedTxs // NOTE
	suite.Require().Equal(newParams, params)
}
