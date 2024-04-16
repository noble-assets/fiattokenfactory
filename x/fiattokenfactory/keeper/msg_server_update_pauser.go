package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdatePauser(ctx context.Context, msg *types.MsgUpdatePauser) (*types.MsgUpdatePauserResponse, error) {
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

	pauser := types.Pauser{
		Address: msg.Address,
	}

	k.SetPauser(ctx, pauser)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgUpdatePauserResponse{}, err
}
