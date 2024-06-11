package types_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateMsgRemoveMinterController(t *testing.T) {

	validAddress := utils.AccAddress().String()

	testCases := map[string]struct {
		from       string
		controller string
		err        error
	}{
		"happy path": {
			from:       validAddress,
			controller: validAddress,
		},
		"invalid from": {
			from: "invalid address",
			err:  sdkerrors.ErrInvalidAddress,
		},
		"invalid conrtoller": {
			from:       validAddress,
			controller: "invalid address",
			err:        sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockMsg := types.MsgRemoveMinterController{
				From:       tc.from,
				Controller: tc.controller,
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
