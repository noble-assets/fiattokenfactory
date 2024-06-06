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

func TestValidateMsgBurn(t *testing.T) {

	validAddress := utils.AccAddress().String()

	ONE := int64(1)
	nONE := int64(-1)
	ZERO := int64(0)

	testCases := map[string]struct {
		from   string
		amount *int64
		err    error
	}{
		"happy path": {
			from:   validAddress,
			amount: &ONE,
		},
		"invalid from address": {
			from: "invalid address",
			err:  sdkerrors.ErrInvalidAddress,
		},
		"nil amount": {
			from:   validAddress,
			amount: nil,
			err:    sdkerrors.ErrInvalidCoins,
		},
		"negative amount": {
			from:   validAddress,
			amount: &nONE,
			err:    sdkerrors.ErrInvalidCoins,
		},
		"zero amount": {
			from:   validAddress,
			amount: &ZERO,
			err:    sdkerrors.ErrInvalidCoins,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var mockCoin sdk.Coin
			if tc.amount != nil {
				mockCoin = sdk.Coin{
					Amount: math.NewInt(*tc.amount),
				}
			}
			mockBurn := types.MsgBurn{
				From:   tc.from,
				Amount: mockCoin,
			}

			err := mockBurn.ValidateBasic()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
