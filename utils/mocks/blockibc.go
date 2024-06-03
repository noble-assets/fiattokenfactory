package mocks

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/circlefin/noble-fiattokenfactory/x/blockibc"
	fiattokenfactorykeeper "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	fiattokenfactorytypes "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	portkeeper "github.com/cosmos/ibc-go/v8/modules/core/05-port/keeper"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/cosmos/ibc-go/v8/testing/mock"
)

func BlockIBC() (blockibc.IBCMiddleware, *fiattokenfactorykeeper.Keeper, sdk.Context) {
	keys := storetypes.NewKVStoreKeys(capabilitytypes.StoreKey, fiattokenfactorytypes.StoreKey)
	mkeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	ctx := testutil.DefaultContextWithKeys(keys, nil, mkeys)

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	capabilityKeeper := capabilitykeeper.NewKeeper(
		cdc, keys[capabilitytypes.StoreKey], mkeys[capabilitytypes.MemStoreKey],
	)
	portKeeper := portkeeper.NewKeeper(
		capabilityKeeper.ScopeToModule(exported.ModuleName),
	)

	transferAppModule := mock.NewAppModule(&portKeeper)
	transferIBCModule := mock.NewIBCModule(
		&transferAppModule,
		mock.NewIBCApp(
			transfertypes.ModuleName,
			capabilityKeeper.ScopeToModule(transfertypes.ModuleName),
		),
	)

	// override the mock ibc_module OnRecvPacket method since it expects specific packet data to return a successful acknowledgment.
	transferIBCModule.IBCApp.OnRecvPacket = func(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
		return mock.MockAcknowledgement
	}

	ftfKeeper := fiattokenfactorykeeper.NewKeeper(
		cdc, nil, runtime.NewKVStoreService(keys[fiattokenfactorytypes.StoreKey]), MockBankKeeper{},
	)
	ftfKeeper.SetMintingDenom(ctx, fiattokenfactorytypes.MintingDenom{Denom: "uusdc"})
	ftfKeeper.SetPaused(ctx, fiattokenfactorytypes.Paused{Paused: false})

	return blockibc.NewIBCMiddleware(
		transferIBCModule,
		ftfKeeper,
	), ftfKeeper, ctx
}
