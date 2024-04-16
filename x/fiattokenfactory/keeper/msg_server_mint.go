package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func (k msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	return k.Keeper.Mint(ctx, msg)
}

func (k Keeper) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	minter, found := k.GetMinters(ctx, msg.From)
	if !found {
		return nil, errors.Wrapf(types.ErrUnauthorized, "you are not a minter")
	}

	_, addressBz, err := bech32.DecodeAndConvert(msg.From)
	if err != nil {
		return nil, err
	}

	_, found = k.GetBlacklisted(ctx, addressBz)
	if found {
		return nil, errors.Wrapf(types.ErrMint, "minter address is blacklisted")
	}

	_, addressBz, err = bech32.DecodeAndConvert(msg.Address)
	if err != nil {
		return nil, err
	}

	_, found = k.GetBlacklisted(ctx, addressBz)
	if found {
		return nil, errors.Wrapf(types.ErrMint, "receiver address is blacklisted")
	}

	mintingDenom := k.GetMintingDenom(ctx)

	if msg.Amount.Denom != mintingDenom.Denom {
		return nil, errors.Wrapf(types.ErrMint, "minting denom is incorrect")
	}

	if minter.Allowance.IsLT(msg.Amount) {
		return nil, errors.Wrapf(types.ErrMint, "minting amount is greater than the allowance")
	}

	paused := k.GetPaused(ctx)

	if paused.Paused {
		return nil, errors.Wrapf(types.ErrMint, "minting is paused")
	}

	minter.Allowance = minter.Allowance.Sub(msg.Amount)

	k.SetMinters(ctx, minter)

	amount := sdk.NewCoins(msg.Amount)

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
		return nil, errors.Wrap(types.ErrMint, err.Error())
	}

	receiver, _ := sdk.AccAddressFromBech32(msg.Address)

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, amount); err != nil {
		return nil, errors.Wrap(types.ErrSendCoinsToAccount, err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = sdkCtx.EventManager().EmitTypedEvent(msg)

	return &types.MsgMintResponse{}, err
}
