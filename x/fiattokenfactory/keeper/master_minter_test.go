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

func createTestMasterMinter(keeper *keeper.Keeper, ctx sdk.Context) types.MasterMinter {
	item := types.MasterMinter{}
	keeper.SetMasterMinter(ctx, item)
	return item
}

func TestMasterMinterGet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	rst, found := keeper.GetMasterMinter(ctx)
	require.Empty(t, rst)
	require.False(t, found)

	item := createTestMasterMinter(keeper, ctx)
	rst, found = keeper.GetMasterMinter(ctx)
	require.True(t, found)
	require.Equal(t,
		utils.Fill(&item),
		utils.Fill(&rst),
	)
}
