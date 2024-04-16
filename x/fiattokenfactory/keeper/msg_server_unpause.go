package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Unpause(ctx context.Context, msg *types.MsgUnpause) (*types.MsgUnpauseResponse, error) {
	pauser, found := k.GetPauser(ctx)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "pauser is not set")
	}

	if pauser.Address != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not the pauser")
	}

	paused := types.Paused{
		Paused: false,
	}

	k.SetPaused(ctx, paused)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgUnpauseResponse{}, err
}
