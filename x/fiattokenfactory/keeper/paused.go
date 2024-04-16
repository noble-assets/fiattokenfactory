package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetPaused set paused in the store
func (k Keeper) SetPaused(ctx context.Context, paused types.Paused) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&paused)
	store.Set(types.KeyPrefix(types.PausedKey), b)
}

// GetPaused returns paused
func (k Keeper) GetPaused(ctx context.Context) (val types.Paused) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyPrefix(types.PausedKey))
	if b == nil {
		panic("Paused state is not set")
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}
