package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetMinters set a specific minters in the store from its index
func (k Keeper) SetMinters(ctx context.Context, minters types.Minters) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MintersKeyPrefix))
	b := k.cdc.MustMarshal(&minters)
	store.Set(types.MintersKey(
		minters.Address,
	), b)
}

// GetMinters returns a minters from its index
func (k Keeper) GetMinters(
	ctx context.Context,
	address string,
) (val types.Minters, found bool) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MintersKeyPrefix))

	b := store.Get(types.MintersKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMinters removes a minters from the store
func (k Keeper) RemoveMinters(
	ctx context.Context,
	address string,
) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MintersKeyPrefix))
	store.Delete(types.MintersKey(
		address,
	))
}

// GetAllMinters returns all minters
func (k Keeper) GetAllMinters(ctx context.Context) (list []types.Minters) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.MintersKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Minters
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
