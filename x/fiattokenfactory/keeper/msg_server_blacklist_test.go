package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/btcutil/bech32"
	"github.com/stretchr/testify/require"
)

func TestBlacklist(t *testing.T) {
	var (
		mockBlacklister  = utils.TestAccount()
		blacklistedUser1 = utils.TestAccount()
		blacklistedUser2 = utils.TestAccount()
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.Blacklist(ctx, &types.MsgBlacklist{From: mockBlacklister.Address})
	require.ErrorIs(t, err, types.ErrUserNotFound)

	ftf.SetBlacklister(ctx, types.Blacklister{Address: mockBlacklister.Address})
	_, err = msgServer.Blacklist(ctx, &types.MsgBlacklist{From: "notTheBlacklister"})
	require.ErrorIs(t, err, types.ErrUnauthorized)

	_, err = msgServer.Blacklist(ctx, &types.MsgBlacklist{From: mockBlacklister.Address, Address: "invalid address"})
	require.ErrorIs(t, err, bech32.ErrInvalidCharacter(32))

	ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: blacklistedUser1.AddressBz})
	_, err = msgServer.Blacklist(ctx, &types.MsgBlacklist{From: mockBlacklister.Address, Address: blacklistedUser1.Address})
	require.ErrorIs(t, err, types.ErrUserBlacklisted)

	res, err := msgServer.Blacklist(ctx, &types.MsgBlacklist{From: mockBlacklister.Address, Address: blacklistedUser2.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgBlacklistResponse{}, res)

}
