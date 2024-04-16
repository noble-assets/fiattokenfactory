package keeper_test

import (
	"strconv"
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNMinterController(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.MinterController {
	items := make([]types.MinterController, n)
	for i := range items {
		items[i].Controller = strconv.Itoa(i)

		keeper.SetMinterController(ctx, items[i])
	}
	return items
}

func TestMinterControllerGet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	items := createNMinterController(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetMinterController(ctx,
			item.Controller,
		)
		require.True(t, found)
		require.Equal(t,
			utils.Fill(&item),
			utils.Fill(&rst),
		)
	}
}

func TestMinterControllerRemove(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	items := createNMinterController(keeper, ctx, 10)
	for _, item := range items {
		keeper.DeleteMinterController(ctx,
			item.Minter,
		)
		_, found := keeper.GetMinterController(ctx,
			item.Minter,
		)
		require.False(t, found)
	}
}

func TestMinterControllerGetAll(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	items := createNMinterController(keeper, ctx, 10)
	require.ElementsMatch(t,
		utils.Fill(items),
		utils.Fill(keeper.GetAllMinterControllers(ctx)),
	)
}
