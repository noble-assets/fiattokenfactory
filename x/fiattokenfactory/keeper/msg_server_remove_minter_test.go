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

func TestRemoveMinter(t *testing.T) {
	var (
		controller   = utils.TestAccount()
		minter       = utils.TestAccount()
		mintingDenom = "uusdc"
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: mintingDenom})
	ftf.SetPaused(ctx, types.Paused{Paused: true})
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.RemoveMinter(ctx, &types.MsgRemoveMinter{})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "minter controller not found")

	otherMinter := "otherminter"
	ftf.SetMinterController(ctx, types.MinterController{Minter: minter.Address, Controller: controller.Address})
	_, err = msgServer.RemoveMinter(ctx, &types.MsgRemoveMinter{From: controller.Address, Address: otherMinter})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, fmt.Sprintf("minter address ≠ minter controller's minter address, (%s≠%s)", otherMinter, minter.Address))

	_, err = msgServer.RemoveMinter(ctx, &types.MsgRemoveMinter{From: controller.Address, Address: minter.Address})
	require.ErrorIs(t, err, types.ErrUserNotFound)
	require.ErrorContains(t, err, "a minter with a given address doesn't exist")

	ftf.SetMinters(ctx, types.Minters{Address: minter.Address})
	res, err := msgServer.RemoveMinter(ctx, &types.MsgRemoveMinter{From: controller.Address, Address: minter.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgRemoveMinterResponse{}, res)

}
