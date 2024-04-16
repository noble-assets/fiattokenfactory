package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Paused(ctx context.Context, req *types.QueryGetPausedRequest) (*types.QueryGetPausedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val := k.GetPaused(ctx)

	return &types.QueryGetPausedResponse{Paused: val}, nil
}
