package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveMinter(ctx context.Context, msg *types.MsgRemoveMinter) (*types.MsgRemoveMinterResponse, error) {
	minterController, found := k.GetMinterController(ctx, msg.From)
	if !found {
		return nil, errors.Wrapf(types.ErrUnauthorized, "minter controller not found")
	}

	if msg.From != minterController.Controller {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not a controller of this minter")
	}

	if msg.Address != minterController.Minter {
		return nil, errors.Wrapf(
			types.ErrUnauthorized,
			"minter address ≠ minter controller's minter address, (%s≠%s)",
			msg.Address, minterController.Minter,
		)
	}

	minter, found := k.GetMinters(ctx, msg.Address)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "a minter with a given address doesn't exist")
	}

	k.RemoveMinters(ctx, minter.Address)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgRemoveMinterResponse{}, err
}
