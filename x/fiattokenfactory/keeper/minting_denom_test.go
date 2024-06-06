package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createTestMintingDenom(keeper *keeper.Keeper, ctx sdk.Context) types.MintingDenom {
	item := types.MintingDenom{
		Denom: "abcd",
	}
	keeper.SetMintingDenom(ctx, item)
	return item
}

func TestMintingDenomSet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	keeper.SetMintingDenom(ctx, types.MintingDenom{})
	require.Panics(t, func() { keeper.SetMintingDenom(ctx, types.MintingDenom{}) })
}

func TestMintingDenomGet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()

	// minting deonom not set, should panic
	require.Panics(t, func() { keeper.GetMintingDenom(ctx) })

	item := createTestMintingDenom(keeper, ctx)
	rst := keeper.GetMintingDenom(ctx)
	require.Equal(t,
		utils.Fill(&item),
		utils.Fill(&rst),
	)
}
