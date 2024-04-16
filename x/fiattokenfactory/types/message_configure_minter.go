package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgConfigureMinter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
	}

	if msg.Allowance.IsNil() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "allowance amount cannot be nil")
	}

	if msg.Allowance.IsNegative() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "allowance amount cannot be negative")
	}

	return nil
}
