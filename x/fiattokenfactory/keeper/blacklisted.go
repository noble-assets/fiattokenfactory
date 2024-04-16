package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetBlacklisted set a specific blacklisted in the store from its index
func (k Keeper) SetBlacklisted(ctx context.Context, blacklisted types.Blacklisted) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.BlacklistedKeyPrefix))
	b := k.cdc.MustMarshal(&blacklisted)
	store.Set(types.BlacklistedKey(blacklisted.AddressBz), b)
}

// GetBlacklisted returns a blacklisted from its index
func (k Keeper) GetBlacklisted(ctx context.Context, addressBz []byte) (val types.Blacklisted, found bool) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.BlacklistedKeyPrefix))

	b := store.Get(types.BlacklistedKey(addressBz))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBlacklisted removes a blacklisted from the store
func (k Keeper) RemoveBlacklisted(ctx context.Context, addressBz []byte) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.BlacklistedKeyPrefix))
	store.Delete(types.BlacklistedKey(addressBz))
}

// GetAllBlacklisted returns all blacklisted
func (k Keeper) GetAllBlacklisted(ctx context.Context) (list []types.Blacklisted) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.BlacklistedKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Blacklisted
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
