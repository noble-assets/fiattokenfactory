package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MinterControllerAll(ctx context.Context, req *types.QueryAllMinterControllerRequest) (*types.QueryAllMinterControllerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var minterControllers []types.MinterController

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	minterControllerStore := prefix.NewStore(store, types.KeyPrefix(types.MinterControllerKeyPrefix))

	pageRes, err := query.Paginate(minterControllerStore, req.Pagination, func(key []byte, value []byte) error {
		var minterController types.MinterController
		if err := k.cdc.Unmarshal(value, &minterController); err != nil {
			return err
		}

		minterControllers = append(minterControllers, minterController)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMinterControllerResponse{MinterController: minterControllers, Pagination: pageRes}, nil
}

func (k Keeper) MinterController(ctx context.Context, req *types.QueryGetMinterControllerRequest) (*types.QueryGetMinterControllerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetMinterController(
		ctx,
		req.ControllerAddress,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetMinterControllerResponse{MinterController: val}, nil
}
