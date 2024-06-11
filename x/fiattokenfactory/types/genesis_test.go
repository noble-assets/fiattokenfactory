package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateGenesis(t *testing.T) {

	mockAccount := utils.TestAccount()
	mockAllowance := sdk.Coin{
		Denom:  "abcd",
		Amount: math.OneInt(),
	}

	testCases := map[string]struct {
		g                 *types.GenesisState
		expectedErrString string
		expectedErr       error
	}{
		"happy path": {
			g: types.DefaultGenesis(),
		},
		"duplicate blacklisted user": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.BlacklistedList = []types.Blacklisted{
					{
						AddressBz: mockAccount.AddressBz,
					},
					{
						AddressBz: mockAccount.AddressBz,
					},
				}
				return g
			}(),
			expectedErrString: "duplicated index for blacklisted",
		},
		"duplicate minters": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MintersList = []types.Minters{
					{
						Address:   mockAccount.Address,
						Allowance: mockAllowance,
					},
					{
						Address:   mockAccount.Address,
						Allowance: mockAllowance,
					},
				}
				return g
			}(),
			expectedErrString: "duplicated index for minters",
		},
		"invalid minter": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MintersList = []types.Minters{
					{
						Address:   "invalid address",
						Allowance: mockAllowance,
					},
				}
				return g
			}(),
			expectedErr:       sdkerrors.ErrInvalidAddress,
			expectedErrString: "invalid minter address",
		},
		"minter has nil allowance": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MintersList = []types.Minters{
					{
						Address:   mockAccount.Address,
						Allowance: sdk.Coin{Denom: "abcd"},
					},
				}
				return g
			}(),
			expectedErr:       sdkerrors.ErrInvalidCoins,
			expectedErrString: "minter allowance cannot be nil or negative",
		},
		"duplicate minter controllers": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MinterControllerList = []types.MinterController{
					{
						Minter:     mockAccount.Address,
						Controller: mockAccount.Address,
					},
					{
						Minter:     mockAccount.Address,
						Controller: mockAccount.Address,
					},
				}
				return g
			}(),
			expectedErrString: "duplicated index for minterController",
		},
		"minter controller has invalid minter": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MinterControllerList = []types.MinterController{
					{
						Minter: "invalid address",
					},
				}
				return g
			}(),
			expectedErr:       sdkerrors.ErrInvalidAddress,
			expectedErrString: "minter controller has invalid minter address",
		},
		"minter controller has invalid controller": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MinterControllerList = []types.MinterController{
					{
						Minter:     mockAccount.Address,
						Controller: "invalid address",
					},
				}
				return g
			}(),
			expectedErr:       sdkerrors.ErrInvalidAddress,
			expectedErrString: "minter controller has invalid controller address",
		},
		"valid owner": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Owner = &types.Owner{
					Address: mockAccount.Address,
				}
				return g
			}(),
		},
		"invalid owner": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Owner = &types.Owner{
					Address: "invalid address",
				}
				return g
			}(),
			expectedErrString: "invalid owner address",
			expectedErr:       sdkerrors.ErrInvalidAddress,
		},
		"valid master minter": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MasterMinter = &types.MasterMinter{
					Address: mockAccount.Address,
				}
				return g
			}(),
		},
		"invalid master minter": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MasterMinter = &types.MasterMinter{
					Address: "invalid address",
				}
				return g
			}(),
			expectedErrString: "invalid master minter address",
			expectedErr:       sdkerrors.ErrInvalidAddress,
		},
		"valid pauser": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Pauser = &types.Pauser{
					Address: mockAccount.Address,
				}
				return g
			}(),
		},
		"invalid pauser": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Pauser = &types.Pauser{
					Address: "invalid address",
				}
				return g
			}(),
			expectedErrString: "invalid pauser address",
			expectedErr:       sdkerrors.ErrInvalidAddress,
		},
		"valid blacklister": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Blacklister = &types.Blacklister{
					Address: mockAccount.Address,
				}
				return g
			}(),
		},
		"invalid blacklister": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Blacklister = &types.Blacklister{
					Address: "invalid address",
				}
				return g
			}(),
			expectedErrString: "invalid black lister address",
			expectedErr:       sdkerrors.ErrInvalidAddress,
		},
		"account with multiple privledges": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.Owner = &types.Owner{
					Address: mockAccount.Address,
				}
				g.MasterMinter = &types.MasterMinter{
					Address: mockAccount.Address,
				}
				return g
			}(),
			expectedErr: types.ErrAlreadyPrivileged,
		},
		"empty minting denom": {
			g: func() *types.GenesisState {
				g := types.DefaultGenesis()
				g.MintingDenom = &types.MintingDenom{
					Denom: "",
				}
				return g
			}(),
			expectedErrString: "minting denom cannot be an empty string",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.g.Validate()

			if tc.expectedErrString != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedErrString)
			} else if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
