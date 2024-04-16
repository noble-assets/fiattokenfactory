package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetPauser set pauser in the store
func (k Keeper) SetPauser(ctx context.Context, pauser types.Pauser) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&pauser)
	store.Set(types.KeyPrefix(types.PauserKey), b)
}

// GetPauser returns pauser
func (k Keeper) GetPauser(ctx context.Context) (val types.Pauser, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyPrefix(types.PauserKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
