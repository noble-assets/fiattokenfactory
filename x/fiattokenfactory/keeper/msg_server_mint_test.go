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

func TestMint(t *testing.T) {
	var (
		receiver     = utils.TestAccount()
		minter       = utils.TestAccount()
		mintingDenom = "uusdc"
		allowance    = sdk.Coin{Denom: mintingDenom, Amount: math.NewInt(10)}
		sendAmount   = sdk.Coin{Denom: mintingDenom, Amount: math.OneInt()}
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: mintingDenom})
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.Mint(ctx, &types.MsgMint{})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "you are not a minter")

	ftf.SetMinters(ctx, types.Minters{Address: "invalid address"})
	_, err = msgServer.Mint(ctx, &types.MsgMint{From: "invalid address"})
	require.Error(t, err)

	ftf.SetMinters(ctx, types.Minters{Address: minter.Address, Allowance: allowance})
	ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: minter.AddressBz})
	_, err = msgServer.Mint(ctx, &types.MsgMint{From: minter.Address})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "minter address is blacklisted")

	ftf.RemoveBlacklisted(ctx, minter.AddressBz)
	_, err = msgServer.Mint(ctx, &types.MsgMint{From: minter.Address, Address: "invalid address"})
	require.Error(t, err)

	ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: receiver.AddressBz})
	_, err = msgServer.Mint(ctx, &types.MsgMint{From: minter.Address, Address: receiver.Address})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "receiver address is blacklisted")

	ftf.RemoveBlacklisted(ctx, receiver.AddressBz)
	_, err = msgServer.Mint(ctx, &types.MsgMint{
		From:    minter.Address,
		Address: receiver.Address,
		Amount:  sdk.Coin{Denom: "notMintingDenom"}})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "minting denom is incorrect")

	_, err = msgServer.Mint(ctx, &types.MsgMint{
		From:    minter.Address,
		Address: receiver.Address,
		Amount:  sdk.Coin{Denom: mintingDenom, Amount: math.NewInt(20)}})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "minting amount is greater than the allowance")

	ftf.SetPaused(ctx, types.Paused{Paused: true})
	_, err = msgServer.Mint(ctx, &types.MsgMint{
		From:    minter.Address,
		Address: receiver.Address,
		Amount:  sendAmount,
	})
	require.ErrorIs(t, err, types.ErrMint)
	require.ErrorContains(t, err, "minting is paused")

	ftf.SetPaused(ctx, types.Paused{Paused: false})
	res, err := msgServer.Mint(ctx, &types.MsgMint{
		From:    minter.Address,
		Address: receiver.Address,
		Amount:  sendAmount,
	})
	require.NoError(t, err)
	require.Equal(t, &types.MsgMintResponse{}, res)

}
