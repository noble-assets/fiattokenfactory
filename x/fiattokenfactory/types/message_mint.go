package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "mint amount cannot be nil")
	}

	if msg.Amount.IsNegative() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "mint amount cannot be negative")
	}

	if msg.Amount.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "mint amount cannot be zero")
	}

	return nil
}
