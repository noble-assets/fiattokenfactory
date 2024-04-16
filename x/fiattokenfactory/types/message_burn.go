package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "burn amount cannot be nil")
	}

	if msg.Amount.IsNegative() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "burn amount cannot be negative")
	}

	if msg.Amount.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "burn amount cannot be zero")
	}

	return nil
}
