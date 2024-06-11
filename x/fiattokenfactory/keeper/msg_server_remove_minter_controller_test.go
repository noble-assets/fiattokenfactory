package keeper_test

import (
	"fmt"
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestRemoveMinterController(t *testing.T) {
	var masterMinter = utils.TestAccount()
	var controller = utils.TestAccount()
	var minter = utils.TestAccount()

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.RemoveMinterController(ctx, &types.MsgRemoveMinterController{})
	require.ErrorIs(t, types.ErrUserNotFound, err)
	require.ErrorContains(t, err, "master minter is not set")

	ftf.SetMasterMinter(ctx, types.MasterMinter{Address: masterMinter.Address})
	_, err = msgServer.RemoveMinterController(ctx, &types.MsgRemoveMinterController{From: "notmasterminter"})
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.ErrorContains(t, err, "you are not the master minter")

	dne := "doesnotexist"
	_, err = msgServer.RemoveMinterController(ctx, &types.MsgRemoveMinterController{From: masterMinter.Address, Controller: dne})
	require.ErrorIs(t, types.ErrUserNotFound, err)
	require.ErrorContains(t, err, fmt.Sprintf("minter controller with a given address (%s) doesn't exist", dne))

	ftf.SetMinterController(ctx, types.MinterController{Minter: minter.Address, Controller: controller.Address})
	res, err := msgServer.RemoveMinterController(ctx, &types.MsgRemoveMinterController{From: masterMinter.Address, Controller: controller.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgRemoveMinterControllerResponse{}, res)

}
