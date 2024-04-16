package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MintingDenom(ctx context.Context, req *types.QueryGetMintingDenomRequest) (*types.QueryGetMintingDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val := k.GetMintingDenom(ctx)

	return &types.QueryGetMintingDenomResponse{MintingDenom: val}, nil
}
