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

func createTestPaused(keeper *keeper.Keeper, ctx sdk.Context) types.Paused {
	item := types.Paused{}
	keeper.SetPaused(ctx, item)
	return item
}

func TestPausedGet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()

	// keeper not set, should panic
	require.Panics(t, func() { keeper.GetPaused(ctx) })

	item := createTestPaused(keeper, ctx)
	rst := keeper.GetPaused(ctx)
	require.Equal(t,
		utils.Fill(&item),
		utils.Fill(&rst),
	)
}
