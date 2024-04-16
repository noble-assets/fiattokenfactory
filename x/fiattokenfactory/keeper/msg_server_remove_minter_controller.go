package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
)

func (k msgServer) RemoveMinterController(ctx context.Context, msg *types.MsgRemoveMinterController) (*types.MsgRemoveMinterControllerResponse, error) {
	masterMinter, found := k.GetMasterMinter(ctx)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "master minter is not set")
	}

	if msg.From != masterMinter.Address {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not the master minter")
	}

	_, found = k.GetMinterController(ctx, msg.Controller)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "minter controller with a given address (%s) doesn't exist", msg.Controller)
	}

	k.DeleteMinterController(ctx, msg.Controller)

	return &types.MsgRemoveMinterControllerResponse{}, nil
}
