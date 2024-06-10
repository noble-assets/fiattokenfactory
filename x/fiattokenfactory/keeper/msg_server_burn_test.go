package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"
)

func TestBurn(t *testing.T) {
	var (
		minter       = utils.TestAccount()
		mintingDenom = "uusdc"
		amount       = sdk.Coin{
			Denom:  mintingDenom,
			Amount: math.OneInt(),
		}
	)
	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: mintingDenom})
	ftf.SetPaused(ctx, types.Paused{Paused: false})
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.Burn(ctx, &types.MsgBurn{From: "notMinter"})
	require.ErrorIs(t, err, types.ErrBurn)
	require.ErrorContains(t, err, "you are not a minter")

	ftf.SetMinters(ctx, types.Minters{Address: "invalid address"})
	_, err = msgServer.Burn(ctx, &types.MsgBurn{From: "invalid address"})
	require.Error(t, err)

	ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: minter.AddressBz})
	ftf.SetMinters(ctx, types.Minters{Address: minter.Address})
	_, err = msgServer.Burn(ctx, &types.MsgBurn{From: minter.Address, Amount: amount})
	require.ErrorIs(t, err, types.ErrBurn)
	require.ErrorContains(t, err, "minter address is blacklisted")

	ftf.RemoveBlacklisted(ctx, minter.AddressBz)
	_, err = msgServer.Burn(ctx, &types.MsgBurn{From: minter.Address, Amount: sdk.Coin{Denom: "notDenom"}})
	require.ErrorIs(t, err, types.ErrBurn)
	require.ErrorContains(t, err, "burning denom is incorrect")

	ftf.SetPaused(ctx, types.Paused{Paused: true})
	_, err = msgServer.Burn(ctx, &types.MsgBurn{From: minter.Address, Amount: amount})
	require.ErrorIs(t, err, types.ErrBurn)
	require.ErrorContains(t, err, "burning is paused")

	// create account that cannot be created as an AccAddress from a Bech32 string.
	invalidMinter := make([]byte, 256)
	address := sdk.AccAddress(invalidMinter).String()
	_, bz, _ := bech32.DecodeAndConvert(address)
	invalidAccount := utils.Account{
		Address:   address,
		AddressBz: bz,
	}
	ftf.SetPaused(ctx, types.Paused{Paused: false})
	ftf.SetMinters(ctx, types.Minters{Address: invalidAccount.Address})
	_, err = msgServer.Burn(ctx, &types.MsgBurn{From: invalidAccount.Address, Amount: amount})
	require.ErrorIs(t, err, types.ErrBurn)

	res, err := msgServer.Burn(ctx, &types.MsgBurn{From: minter.Address, Amount: amount})
	require.NoError(t, err)
	require.Equal(t, &types.MsgBurnResponse{}, res)

}
