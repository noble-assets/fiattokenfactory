package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgConfigureMinterController) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Controller)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid controller address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
	}
	return nil
}
