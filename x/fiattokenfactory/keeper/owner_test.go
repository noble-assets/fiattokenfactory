package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

func TestOwnerGet(t *testing.T) {
	keeper, ctx := mocks.FiatTokenfactoryKeeper()

	rst, found := keeper.GetOwner(ctx)
	require.False(t, found)
	require.Empty(t, rst)

	owner := types.Owner{Address: "1"}
	keeper.SetOwner(ctx, owner)

	rst, found = keeper.GetOwner(ctx)
	require.True(t, found)
	require.Equal(t,
		owner,
		utils.Fill(&rst),
	)

	rst, found = keeper.GetPendingOwner(ctx)
	require.False(t, found)
	require.Empty(t, rst)

	newOwner := types.Owner{Address: "2"}

	keeper.SetPendingOwner(ctx, newOwner)

	rst, found = keeper.GetPendingOwner(ctx)
	require.True(t, found)
	require.Equal(t,
		newOwner,
		utils.Fill(&rst),
	)
}
