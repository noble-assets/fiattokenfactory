package types_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateMsgAcceptOwner(t *testing.T) {

	testCases := map[string]struct {
		address string
		err     error
	}{
		"happy path": {
			address: utils.AccAddress().String(),
		},
		"invalid address": {
			address: "invalid address",
			err:     sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockMsg := types.MsgAcceptOwner{From: tc.address}
			err := mockMsg.ValidateBasic()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
