package mocks

import (
	"context"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type MockBankKeeper struct{}

var _ types.BankKeeper = MockBankKeeper{}

func (m MockBankKeeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	// TODO implement me
	panic("implement me")
}

func (m MockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	// TODO implement me
	panic("implement me")
}

func (m MockBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	// TODO implement me
	panic("implement me")
}

func (m MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	// TODO implement me
	panic("implement me")
}

func (m MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	// TODO implement me
	panic("implement me")
}

func (m MockBankKeeper) GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool) {
	return banktypes.Metadata{}, true
}
