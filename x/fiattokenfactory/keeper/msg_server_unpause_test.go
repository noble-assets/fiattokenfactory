package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestUnPause(t *testing.T) {
	var pauser = utils.TestAccount()

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.Unpause(ctx, &types.MsgUnpause{})
	require.ErrorIs(t, err, types.ErrUserNotFound)
	require.ErrorContains(t, err, "pauser is not set")

	ftf.SetPauser(ctx, types.Pauser{Address: pauser.Address})
	_, err = msgServer.Unpause(ctx, &types.MsgUnpause{From: "notpauser"})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "you are not the pauser")

	res, err := msgServer.Unpause(ctx, &types.MsgUnpause{From: pauser.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgUnpauseResponse{}, res)
}
