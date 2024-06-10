package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestConfigureMinter(t *testing.T) {
	var (
		controller   = utils.TestAccount()
		minter       = utils.TestAccount()
		mintingDenom = "uusdc"
		allowance    = sdk.Coin{Denom: mintingDenom, Amount: math.OneInt()}
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: "uusdc"})
	ftf.SetPaused(ctx, types.Paused{Paused: true})
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.ConfigureMinter(ctx, &types.MsgConfigureMinter{})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "minting denom is incorrect")

	_, err = msgServer.ConfigureMinter(ctx, &types.MsgConfigureMinter{From: "notConfiguredMinter", Allowance: allowance})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "minter controller not found")

	ftf.SetMinterController(ctx, types.MinterController{Controller: controller.Address, Minter: minter.Address})
	_, err = msgServer.ConfigureMinter(ctx, &types.MsgConfigureMinter{From: controller.Address, Allowance: allowance})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "minting is paused")

	ftf.SetPaused(ctx, types.Paused{Paused: false})
	_, err = msgServer.ConfigureMinter(ctx, &types.MsgConfigureMinter{From: controller.Address, Address: "unControlledAddress", Allowance: allowance})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "minter address ≠ minter controller's minter address")

	res, err := msgServer.ConfigureMinter(ctx, &types.MsgConfigureMinter{From: controller.Address, Address: minter.Address, Allowance: allowance})
	require.NoError(t, err)
	require.Equal(t, &types.MsgConfigureMinterResponse{}, res)
}
