package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
)

func (k msgServer) ConfigureMinterController(ctx context.Context, msg *types.MsgConfigureMinterController) (*types.MsgConfigureMinterControllerResponse, error) {
	masterMinter, found := k.GetMasterMinter(ctx)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "master minter is not set")
	}

	if masterMinter.Address != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not the master minter")
	}

	controller := types.MinterController{
		Minter:     msg.Minter,
		Controller: msg.Controller,
	}

	k.SetMinterController(ctx, controller)

	return &types.MsgConfigureMinterControllerResponse{}, nil
}
