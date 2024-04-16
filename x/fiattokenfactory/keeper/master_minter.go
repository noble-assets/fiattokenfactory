package keeper

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetMasterMinter set masterMinter in the store
func (k Keeper) SetMasterMinter(ctx context.Context, masterMinter types.MasterMinter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&masterMinter)
	store.Set(types.KeyPrefix(types.MasterMinterKey), b)
}

// GetMasterMinter returns masterMinter
func (k Keeper) GetMasterMinter(ctx context.Context) (val types.MasterMinter, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.KeyPrefix(types.MasterMinterKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
