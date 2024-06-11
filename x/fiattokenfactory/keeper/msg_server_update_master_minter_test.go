package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateMasterMinter(t *testing.T) {
	var (
		owner        = utils.TestAccount()
		masterMinter = utils.TestAccount()
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.UpdateMasterMinter(ctx, &types.MsgUpdateMasterMinter{})
	require.ErrorIs(t, err, types.ErrUserNotFound)
	require.ErrorContains(t, err, "owner is not set")

	ftf.SetOwner(ctx, types.Owner{Address: owner.Address})
	_, err = msgServer.UpdateMasterMinter(ctx, &types.MsgUpdateMasterMinter{From: "nottheowner"})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "you are not the owner")

	_, err = msgServer.UpdateMasterMinter(ctx, &types.MsgUpdateMasterMinter{From: owner.Address, Address: owner.Address})
	require.ErrorIs(t, err, types.ErrAlreadyPrivileged)

	res, err := msgServer.UpdateMasterMinter(ctx, &types.MsgUpdateMasterMinter{From: owner.Address, Address: masterMinter.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgUpdateMasterMinterResponse{}, res)

}
