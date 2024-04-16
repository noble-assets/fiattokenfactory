package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (k msgServer) Unblacklist(ctx context.Context, msg *types.MsgUnblacklist) (*types.MsgUnblacklistResponse, error) {
	blacklister, found := k.GetBlacklister(ctx)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "blacklister is not set")
	}

	if blacklister.Address != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not the blacklister")
	}

	_, addressBz, err := bech32.DecodeAndConvert(msg.Address)
	if err != nil {
		return nil, err
	}

	blacklisted, found := k.GetBlacklisted(ctx, addressBz)
	if !found {
		return nil, errors.Wrapf(types.ErrUserNotFound, "the specified address is not blacklisted")
	}

	k.RemoveBlacklisted(ctx, blacklisted.AddressBz)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgUnblacklistResponse{}, err
}
