package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AcceptOwner(ctx context.Context, msg *types.MsgAcceptOwner) (*types.MsgAcceptOwnerResponse, error) {
	owner, found := k.GetPendingOwner(ctx)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "pending owner is not set")
	}

	if owner.Address != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not the pending owner")
	}

	k.SetOwner(ctx, owner)

	k.DeletePendingOwner(ctx)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgAcceptOwnerResponse{}, err
}
