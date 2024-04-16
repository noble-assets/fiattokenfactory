package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/fiattokenfactory module sentinel errors
var (
	ErrUnauthorized       = errors.Register(ModuleName, 2, "unauthorized")
	ErrUserNotFound       = errors.Register(ModuleName, 3, "user not found")
	ErrMint               = errors.Register(ModuleName, 4, "tokens can not be minted")
	ErrSendCoinsToAccount = errors.Register(ModuleName, 5, "can't send tokens to account")
	ErrBurn               = errors.Register(ModuleName, 6, "tokens can not be burned")
	ErrPaused             = errors.Register(ModuleName, 7, "the chain is paused")
	ErrMintingDenomSet    = errors.Register(ModuleName, 9, "the minting denom has already been set")
	ErrUserBlacklisted    = errors.Register(ModuleName, 10, "user is already blacklisted")
	ErrAlreadyPrivileged  = errors.Register(ModuleName, 11, "address is already assigned to privileged role")
	ErrDenomNotRegistered = errors.Register(ModuleName, 12, "denom not registered in bank module")
)
