package fiattokenfactory_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	t.Parallel()

	mockAccount := utils.TestAccount()
	mockAllowance := sdk.Coin{
		Denom:  "abcd",
		Amount: math.OneInt(),
	}

	ftfKeeper, ctx := mocks.FiatTokenfactoryKeeper()
	bk := mocks.MockBankKeeper{}
	g := types.GenesisState{
		BlacklistedList: []types.Blacklisted{
			{
				AddressBz: mockAccount.AddressBz,
			},
		},
		Paused:       &types.Paused{Paused: false},
		MasterMinter: &types.MasterMinter{Address: mockAccount.Address},
		MintersList: []types.Minters{
			{
				Address:   mockAccount.Address,
				Allowance: mockAllowance,
			},
		},
		Pauser:      &types.Pauser{Address: mockAccount.Address},
		Blacklister: &types.Blacklister{Address: mockAccount.Address},
		Owner:       &types.Owner{Address: mockAccount.Address},
		MinterControllerList: []types.MinterController{
			{
				Minter:     mockAccount.Address,
				Controller: mockAccount.Address,
			},
		},
		MintingDenom: &types.MintingDenom{
			Denom: "uusdc",
		},
	}

	require.NotPanics(t, func() { fiattokenfactory.InitGenesis(ctx, ftfKeeper, bk, g) })

	// set minting denom not set in bank module
	g.MintingDenom = &types.MintingDenom{Denom: "abcd"}

	require.Panics(t, func() { fiattokenfactory.InitGenesis(ctx, ftfKeeper, bk, g) })

}

func TestExportGenesis(t *testing.T) {
	t.Parallel()
	ftfKeeper, ctx := mocks.FiatTokenfactoryKeeper()
	ftfKeeper.SetMintingDenom(ctx, types.MintingDenom{Denom: "uusdc"})
	ftfKeeper.SetPaused(ctx, types.Paused{Paused: false})
	ftfKeeper.SetMasterMinter(ctx, types.MasterMinter{Address: "mock"})
	ftfKeeper.SetPauser(ctx, types.Pauser{Address: "mock"})
	ftfKeeper.SetBlacklister(ctx, types.Blacklister{Address: "mock"})
	ftfKeeper.SetOwner(ctx, types.Owner{Address: "mock"})

	require.NotPanics(t, func() { fiattokenfactory.ExportGenesis(ctx, ftfKeeper) })
}
