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

func createTestPauser(keeper *keeper.Keeper, ctx sdk.Context) types.Pauser {
	item := types.Pauser{}
	keeper.SetPauser(ctx, item)
	return item
}

func TestPauserGet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()
	item := createTestPauser(keeper, ctx)
	rst, found := keeper.GetPauser(ctx)
	require.True(t, found)
	require.Equal(t,
		utils.Fill(&item),
		utils.Fill(&rst),
	)
}
