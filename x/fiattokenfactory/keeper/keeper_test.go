package keeper_test

import (
	"testing"

	cosmossdk_io_math "cosmossdk.io/math"
	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeperLogger(t *testing.T) {
	t.Parallel()
	k, _ := mocks.FiatTokenfactoryKeeper()
	l := k.Logger()
	require.NotNil(t, l)
}

func TestSendRestrictionsFn(t *testing.T) {

	// ARRANGE: Mock sender and receiver.
	from, to := utils.AccAddress(), utils.AccAddress()

	// ARRANGE: Organize table driven test cases.
	testCases := map[string]struct {
		tfPaused         bool
		fromAddr, toAddr sdk.AccAddress
		toBlacklist      *sdk.AccAddress
		expectedErr      error
	}{
		"happy path": {
			tfPaused:    false,
			fromAddr:    from,
			toAddr:      to,
			toBlacklist: nil,
			expectedErr: nil,
		},
		"tokenfactory paused": {
			tfPaused:    true,
			fromAddr:    from,
			toAddr:      to,
			toBlacklist: nil,
			expectedErr: types.ErrPaused,
		},
		"sender blacklisted": {
			tfPaused:    false,
			fromAddr:    from,
			toAddr:      to,
			toBlacklist: &from,
			expectedErr: types.ErrUnauthorized,
		},
		"receiver blacklisted": {
			tfPaused:    false,
			fromAddr:    from,
			toAddr:      to,
			toBlacklist: &to,
			expectedErr: types.ErrUnauthorized,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			// ARRANGE: Mock fiattokenfactory.
			k, ctx := mocks.FiatTokenfactoryKeeper()
			k.SetMintingDenom(ctx, types.MintingDenom{Denom: "uusdc"})

			// ACT: Set paused and blacklisted state based on test case.
			k.SetPaused(ctx, types.Paused{Paused: tc.tfPaused})
			if tc.toBlacklist != nil {
				k.SetBlacklisted(ctx, types.Blacklisted{
					AddressBz: *tc.toBlacklist,
				})
			}

			//ACT: Test send restrictions
			amount := sdk.NewCoins(sdk.NewCoin(k.GetMintingDenom(ctx).Denom, cosmossdk_io_math.OneInt()))
			_, err := k.SendRestrictionFn(ctx, tc.fromAddr, tc.toAddr, amount)

			// Assert: Assert proper error based on test case
			if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})

	}
}

func TestValidatePrivileges(t *testing.T) {

	// ARRANGE: Mock fiattokenfactory.
	k, ctx := mocks.FiatTokenfactoryKeeper()

	// ARRANGE: Create wallets.
	owner := utils.TestAccount()
	blacklister := utils.TestAccount()
	masterMinter := utils.TestAccount()
	pauser := utils.TestAccount()

	address := []utils.Account{owner, blacklister, masterMinter, pauser}

	// ARRANGE: Set tokenfacotry privledges
	k.SetOwner(ctx, types.Owner{Address: owner.Address})
	k.SetBlacklister(ctx, types.Blacklister{Address: blacklister.Address})
	k.SetMasterMinter(ctx, types.MasterMinter{Address: masterMinter.Address})
	k.SetPauser(ctx, types.Pauser{Address: pauser.Address})

	// ACT: Assert error on already privledged accounts
	for _, ad := range address {
		err := k.ValidatePrivileges(ctx, ad.Address)
		require.ErrorIs(t, err, types.ErrAlreadyPrivileged)
	}

	// ACT: Assert error when passing malformed address
	err := k.ValidatePrivileges(ctx, "malformed bech32 address")
	require.Error(t, err)

	// ARRANGE: Create new account with no privledges
	alice := utils.TestAccount()

	//ACT: Assert no error for non privledged account
	err = k.ValidatePrivileges(ctx, alice.Address)
	require.NoError(t, err)
}
