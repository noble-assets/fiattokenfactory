package types_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateMsgUpdateBlacklister(t *testing.T) {

	validAddress := utils.AccAddress().String()

	testCases := map[string]struct {
		from    string
		address string
		err     error
	}{
		"happy path": {
			from:    validAddress,
			address: validAddress,
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
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockMsg := types.MsgUpdateBlacklister{
				From:    tc.from,
				Address: tc.address,
			}
			err := mockMsg.ValidateBasic()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
