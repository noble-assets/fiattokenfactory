package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestUnblacklist(t *testing.T) {
	var (
		blacklister = utils.TestAccount()
		blUser      = utils.TestAccount()
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: "uusdc"})
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.Unblacklist(ctx, &types.MsgUnblacklist{})
	require.ErrorIs(t, err, types.ErrUserNotFound)
	require.ErrorContains(t, err, "blacklister is not set")

	ftf.SetBlacklister(ctx, types.Blacklister{Address: blacklister.Address})
	_, err = msgServer.Unblacklist(ctx, &types.MsgUnblacklist{From: "notTheBlacklister"})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "you are not the blacklister")

	_, err = msgServer.Unblacklist(ctx, &types.MsgUnblacklist{From: blacklister.Address, Address: "invalid address"})
	require.Error(t, err)

	_, err = msgServer.Unblacklist(ctx, &types.MsgUnblacklist{From: blacklister.Address, Address: blUser.Address})
	require.ErrorIs(t, types.ErrUserNotFound, err)
	require.ErrorContains(t, err, "the specified address is not blacklisted")

	ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: blUser.AddressBz})
	res, err := msgServer.Unblacklist(ctx, &types.MsgUnblacklist{From: blacklister.Address, Address: blUser.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgUnblacklistResponse{}, res)
}
