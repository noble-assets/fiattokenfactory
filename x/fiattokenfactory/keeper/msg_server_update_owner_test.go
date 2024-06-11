package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateOwner(t *testing.T) {
	var (
		owner    = utils.TestAccount()
		newOwner = utils.TestAccount()
	)

	ftf, ctx := mocks.FiatTokenfactoryKeeper()
	msgServer := keeper.NewMsgServerImpl(ftf)

	_, err := msgServer.UpdateOwner(ctx, &types.MsgUpdateOwner{})
	require.ErrorIs(t, err, types.ErrUserNotFound)
	require.ErrorContains(t, err, "owner is not set")

	ftf.SetOwner(ctx, types.Owner{Address: owner.Address})
	_, err = msgServer.UpdateOwner(ctx, &types.MsgUpdateOwner{From: "nottheowner"})
	require.ErrorIs(t, err, types.ErrUnauthorized)
	require.ErrorContains(t, err, "you are not the owner")

	_, err = msgServer.UpdateOwner(ctx, &types.MsgUpdateOwner{From: owner.Address, Address: owner.Address})
	require.ErrorIs(t, err, types.ErrAlreadyPrivileged)

	ftf.SetMasterMinter(ctx, types.MasterMinter{Address: newOwner.Address})
	_, err = msgServer.UpdateOwner(ctx, &types.MsgUpdateOwner{From: owner.Address, Address: newOwner.Address})
	require.ErrorIs(t, err, types.ErrAlreadyPrivileged)

	ftf.SetMasterMinter(ctx, types.MasterMinter{Address: "masterminter"})
	res, err := msgServer.UpdateOwner(ctx, &types.MsgUpdateOwner{From: owner.Address, Address: newOwner.Address})
	require.NoError(t, err)
	require.Equal(t, &types.MsgUpdateOwnerResponse{}, res)

}
