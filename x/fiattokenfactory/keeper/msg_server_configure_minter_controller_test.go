package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestConfigureMinterController(t *testing.T) {
	var masterMinter = utils.TestAccount()

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.ConfigureMinterController(ctx, &types.MsgConfigureMinterController{})
	require.ErrorIs(t, err, types.ErrUserNotFound)
	require.ErrorContains(t, err, "master minter is not set")

	ftf.SetMasterMinter(ctx, types.MasterMinter{Address: masterMinter.Address})
	_, err = msgServer.ConfigureMinterController(ctx, &types.MsgConfigureMinterController{From: "notMasterMinter"})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "you are not the master minter")

	_, err = msgServer.ConfigureMinterController(ctx, &types.MsgConfigureMinterController{From: masterMinter.Address})
	require.NoError(t, err)
}
