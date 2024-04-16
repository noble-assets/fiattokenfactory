package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetBlacklister set blacklister in the store
func (k Keeper) SetBlacklister(ctx context.Context, blacklister types.Blacklister) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&blacklister)
	store.Set(types.KeyPrefix(types.BlacklisterKey), b)
}

// GetBlacklister returns blacklister
func (k Keeper) GetBlacklister(ctx context.Context) (val types.Blacklister, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyPrefix(types.BlacklisterKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
