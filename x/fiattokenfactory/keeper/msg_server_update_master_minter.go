package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateMasterMinter(ctx context.Context, msg *types.MsgUpdateMasterMinter) (*types.MsgUpdateMasterMinterResponse, error) {
	owner, found := k.GetOwner(ctx)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "owner is not set")
	}

	if owner.Address != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not the owner")
	}

	// ensure that the specified address is not already assigned to a privileged role
	err := k.ValidatePrivileges(ctx, msg.Address)
	if err != nil {
		return nil, err
	}

	masterMinter := types.MasterMinter{
		Address: msg.Address,
	}

	k.SetMasterMinter(ctx, masterMinter)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgUpdateMasterMinterResponse{}, err
}
