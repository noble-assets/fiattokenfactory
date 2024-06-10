package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestAccpetOwner(t *testing.T) {
	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.AcceptOwner(ctx, &types.MsgAcceptOwner{From: "mock"})
	require.ErrorIs(t, types.ErrUserNotFound, err)

	ftf.SetPendingOwner(ctx, types.Owner{Address: "mock"})
	_, err = msgServer.AcceptOwner(ctx, &types.MsgAcceptOwner{From: "notFrom"})
	require.ErrorIs(t, err, types.ErrUnauthorized)

	ftf.SetPendingOwner(ctx, types.Owner{Address: "mock"})
	res, err := msgServer.AcceptOwner(ctx, &types.MsgAcceptOwner{From: "mock"})
	require.NoError(t, err)
	require.Equal(t, &types.MsgAcceptOwnerResponse{}, res)

}
