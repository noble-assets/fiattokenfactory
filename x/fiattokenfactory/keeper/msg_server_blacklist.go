package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (k msgServer) Blacklist(ctx context.Context, msg *types.MsgBlacklist) (*types.MsgBlacklistResponse, error) {
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

	_, found = k.GetBlacklisted(ctx, addressBz)
	if found {
		return nil, types.ErrUserBlacklisted
	}

	blacklisted := types.Blacklisted{
		AddressBz: addressBz,
	}

	k.SetBlacklisted(ctx, blacklisted)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgBlacklistResponse{}, err
}
