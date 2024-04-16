package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetMinterController set a specific minterController in the store from its index
func (k Keeper) SetMinterController(ctx context.Context, minterController types.MinterController) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MinterControllerKeyPrefix))
	b := k.cdc.MustMarshal(&minterController)
	store.Set(types.MinterControllerKey(
		minterController.Controller,
	), b)
}

// GetMinterController returns a minterController from its index
func (k Keeper) GetMinterController(
	ctx context.Context,
	controller string,
) (val types.MinterController, found bool) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MinterControllerKeyPrefix))

	b := store.Get(types.MinterControllerKey(
		controller,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMinterController removes a minterController from the store
func (k Keeper) DeleteMinterController(
	ctx context.Context,
	controller string,
) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MinterControllerKeyPrefix))
	store.Delete(types.MinterControllerKey(
		controller,
	))
}

// GetAllMinterControllers returns all minterController
func (k Keeper) GetAllMinterControllers(ctx context.Context) (list []types.MinterController) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MinterControllerKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MinterController
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
