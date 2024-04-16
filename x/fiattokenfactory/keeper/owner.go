package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetOwner set owner in the store
func (k Keeper) SetOwner(ctx context.Context, owner types.Owner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&owner)
	store.Set(types.KeyPrefix(types.OwnerKey), b)
}

// GetOwner returns owner
func (k Keeper) GetOwner(ctx context.Context) (val types.Owner, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyPrefix(types.OwnerKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetPendingOwner set pending owner in the store
func (k Keeper) SetPendingOwner(ctx context.Context, owner types.Owner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&owner)
	store.Set(types.KeyPrefix(types.PendingOwnerKey), b)
}

// DeletePendingOwner deletes the pending owner in the store
func (k Keeper) DeletePendingOwner(ctx context.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.KeyPrefix(types.PendingOwnerKey))
}

// GetPendingOwner returns pending owner
func (k Keeper) GetPendingOwner(ctx context.Context) (val types.Owner, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyPrefix(types.PendingOwnerKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
