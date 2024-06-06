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

func TestValidateMsgConfigureMinter(t *testing.T) {

	validAddress := utils.AccAddress().String()

	ONE := int64(1)
	nONE := int64(-1)

	testCases := map[string]struct {
		from      string
		address   string
		allowance *int64
		err       error
	}{
		"happy path": {
			from:      validAddress,
			address:   validAddress,
			allowance: &ONE,
		},
		"invalid from": {
			from: "invalid address",
			err:  sdkerrors.ErrInvalidAddress,
		},
		"invalid address": {
			from:    validAddress,
			address: "invalid address",
			err:     sdkerrors.ErrInvalidAddress,
		},
		"nil allowance": {
			from:      validAddress,
			address:   validAddress,
			allowance: nil,
			err:       sdkerrors.ErrInvalidCoins,
		},
		"negative allowance": {
			from:      validAddress,
			address:   validAddress,
			allowance: &nONE,
			err:       sdkerrors.ErrInvalidCoins,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var mockCoin sdk.Coin
			if tc.allowance != nil {
				mockCoin = sdk.Coin{
					Amount: math.NewInt(*tc.allowance),
				}
			}

			mockConfigureMC := types.MsgConfigureMinter{
				From:      tc.from,
				Address:   tc.address,
				Allowance: mockCoin,
			}

			err := mockConfigureMC.ValidateBasic()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
