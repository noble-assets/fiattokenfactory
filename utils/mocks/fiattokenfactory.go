package mocks

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FiatTokenfactoryKeeper() (*keeper.Keeper, sdk.Context) {
	logger := log.NewNopLogger()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	state := store.NewCommitMultiStore(db.NewMemDB(), logger, metrics.NewNoOpMetrics())
	state.MountStoreWithDB(key, storetypes.StoreTypeIAVL, nil)
	_ = state.LoadLatestVersion()

	return keeper.NewKeeper(
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
		logger,
		runtime.NewKVStoreService(key),
		MockBankKeeper{},
	), sdk.NewContext(state, cmtproto.Header{}, false, logger)
}
