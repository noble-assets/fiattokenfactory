package fiattokenfactory_test

import (
	"time"

	"testing"

	"cosmossdk.io/math"
	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory"
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

var (
	coin                       = sdk.NewInt64Coin("uusdc", 10)
	coins                      = sdk.Coins{coin}
	testAccount1, testAccount2 = utils.TestAccount(), utils.TestAccount()
)

func TestAnteHandlerIsPaused(t *testing.T) {

	// ARRANGE: Arrange table driven test cases
	testCases := map[string]struct {
		expectedFailOnPause bool
		message             sdk.Msg
	}{
		"no message": {
			expectedFailOnPause: false,
		},
		"irrelevant msg": {
			expectedFailOnPause: false,
			message:             &testdata.MsgCreateDog{},
		},
		"msgSend": {
			expectedFailOnPause: true,
			message: &banktypes.MsgSend{
				FromAddress: "mock",
				ToAddress:   "mock",
				Amount:      coins,
			},
		},
		"msgExec": {
			// msgExec itself should not fail on pause, only specific nested messages inside should
			expectedFailOnPause: true,
			message: func() sdk.Msg {
				msgSend := &banktypes.MsgSend{
					FromAddress: "mock",
					ToAddress:   "mock",
					Amount:      coins,
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgSend)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: "mock",
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
		},
		"msgGrant": {
			expectedFailOnPause: true,
			message: func() sdk.Msg {
				mockTime := time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
				mockExpires := mockTime.Add(time.Hour)
				sendAuthz := banktypes.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("uusdc", math.OneInt())), nil)
				sendGrant, err := authz.NewGrant(mockTime, sendAuthz, &mockExpires)
				require.NoError(t, err)

				msg := &authz.MsgGrant{
					Granter: "mock",
					Grantee: "mock",
					Grant:   sendGrant,
				}
				return msg
			}(),
		},
		"msgMultiSend": {
			expectedFailOnPause: true,
			message: func() sdk.Msg {
				multiSendMsg := &banktypes.MsgMultiSend{
					Inputs:  []banktypes.Input{banktypes.NewInput(testAccount1.AddressBz, coins)},
					Outputs: []banktypes.Output{banktypes.NewOutput(testAccount2.AddressBz, coins)},
				}
				return multiSendMsg
			}(),
		},
		"msgTransfer": {
			expectedFailOnPause: true,
			message: &transfertypes.MsgTransfer{
				Token: coin,
			},
		},
		"MsgDepositForBurn": {
			expectedFailOnPause: true,
			message: &cctptypes.MsgDepositForBurn{
				From:      "mock",
				BurnToken: "uusdc",
			},
		},
		"MsgDepositForBurnWithCaller": {
			expectedFailOnPause: true,
			message: &cctptypes.MsgDepositForBurnWithCaller{
				From:      "mock",
				BurnToken: "uusdc",
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// ARRANGE: setup tokenfactory and isPaused decorator
			ftf, ctx := mocks.FiatTokenfactoryKeeper()
			ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: "uusdc"})
			ftf.SetPaused(ctx, types.Paused{Paused: false})
			cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
			ad := fiattokenfactory.NewIsPausedDecorator(cdc, ftf)

			// ARRANGE: Build transactions with specific test case message
			builder, err := newMockTxBuilder(cdc)
			require.NoError(t, err)
			if tc.message != nil {
				err = builder.SetMsgs(tc.message)
				require.NoError(t, err)
			}
			tx := builder.GetTx()

			// ACT: Run transaction through ante handler while chain is NOT paused.
			_, err = ad.AnteHandle(ctx, tx, true, mockNext)

			// ASSERT: No errors while chain is not paused
			require.NoError(t, err)

			// ARRANGE: Pause tokenfactory
			ftf.SetPaused(ctx, types.Paused{Paused: true})

			// ACT: Run transaction through ante handler while chain IS paused.
			_, err = ad.AnteHandle(ctx, tx, true, mockNext)

			// ASSERT: Assert expected isPaused error for specific test case messages
			if tc.expectedFailOnPause {
				require.ErrorIs(t, err, types.ErrPaused)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAnteHandlerIsBlacklisted(t *testing.T) {

	// ARRANGE: Arrange table driven test cases
	testCases := map[string]struct {
		message sdk.Msg
		// if blacklistSendAndRec == true, the test case will run the message through the antehanlder three times:
		// 	- without blacklisting any address
		// 	- blacklisting just the sender
		// 	- blacklisting just the receiver
		blacklistSendAndRec bool
		// set testInvalidAddress to true if testing for an invalid address or an address that
		// cannot be bech32 decoded.
		// If testInvalidAddress is true, blacklistSendAndRec should be false.
		testInvalidAddress bool
		expectedError      error
	}{
		"no message": {
			blacklistSendAndRec: false,
		},
		"irrelevant msg": {
			message:             &testdata.MsgCreateDog{},
			blacklistSendAndRec: false,
		},
		"msgSend": {
			message: &banktypes.MsgSend{
				FromAddress: testAccount1.Address,
				ToAddress:   testAccount2.Address,
				Amount:      coins,
			},
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgSend invalid sender": {
			message: &banktypes.MsgSend{
				FromAddress: testAccount1.Address,
				ToAddress:   "invalid address",
				Amount:      coins,
			},
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgSend invalid receiver": {
			message: &banktypes.MsgSend{
				FromAddress: "invalid address",
				ToAddress:   testAccount2.Address,
				Amount:      coins,
			},
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgMultiSend": {
			message: func() sdk.Msg {
				multiSendMsg := &banktypes.MsgMultiSend{
					Inputs:  []banktypes.Input{banktypes.NewInput(testAccount1.AddressBz, coins)},
					Outputs: []banktypes.Output{banktypes.NewOutput(testAccount2.AddressBz, coins)},
				}
				return multiSendMsg
			}(),
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgMultiSend invalid sender": {
			message: func() sdk.Msg {
				multiSendMsg := &banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: "invalid address",
							Coins:   coins,
						},
					},
					Outputs: []banktypes.Output{banktypes.NewOutput(testAccount2.AddressBz, coins)},
				}
				return multiSendMsg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgMultiSend invalid receiver": {
			message: func() sdk.Msg {
				multiSendMsg := &banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{banktypes.NewInput(testAccount2.AddressBz, coins)},
					Outputs: []banktypes.Output{
						{
							Address: "invalid address",
							Coins:   coins,
						},
					},
				}
				return multiSendMsg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgTransfer": {
			message: &transfertypes.MsgTransfer{
				Sender:   testAccount1.Address,
				Receiver: testAccount2.Address,
				Token:    coin,
			},
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgTransfer invalid sender": {
			message: &transfertypes.MsgTransfer{
				Sender:   "invalid address",
				Receiver: testAccount2.Address,
				Token:    coin,
			},
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgTransfer invalid receiver": {
			message: &transfertypes.MsgTransfer{
				Sender:   testAccount1.Address,
				Receiver: "invalid address",
				Token:    coin,
			},
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgDepositForBurn invalid address": {
			message: &cctptypes.MsgDepositForBurn{
				From:      "invalid address",
				BurnToken: "uusdc",
			},
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgDepositForBurnWithCaller invalid address": {
			message: &cctptypes.MsgDepositForBurnWithCaller{
				From:      "invalid address",
				BurnToken: "uusdc",
			},
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgExec MsgSend": {
			message: func() sdk.Msg {
				mgsSend := &banktypes.MsgSend{
					FromAddress: testAccount2.Address,
					ToAddress:   testAccount2.Address,
					Amount:      coins,
				}
				msgSendAny, err := codectypes.NewAnyWithValue(mgsSend)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: testAccount1.Address,
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgExec MsgSend invalid grantee": {
			message: func() sdk.Msg {
				msgSend := &banktypes.MsgSend{
					FromAddress: testAccount1.Address,
					ToAddress:   testAccount2.Address,
					Amount:      coins,
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgSend)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: "invalid address",
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgExec MultiSend blacklisted grantee": {
			message: func() sdk.Msg {
				multiSendMsg := &banktypes.MsgMultiSend{
					Inputs:  []banktypes.Input{banktypes.NewInput(testAccount1.AddressBz, coins)},
					Outputs: []banktypes.Output{banktypes.NewOutput(testAccount2.AddressBz, coins)},
				}
				msgSendAny, err := codectypes.NewAnyWithValue(multiSendMsg)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: testAccount1.Address,
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgExec MultiSend invalid grantee": {
			message: func() sdk.Msg {
				multiSendMsg := &banktypes.MsgMultiSend{
					Inputs:  []banktypes.Input{banktypes.NewInput(testAccount1.AddressBz, coins)},
					Outputs: []banktypes.Output{banktypes.NewOutput(testAccount2.AddressBz, coins)},
				}
				msgSendAny, err := codectypes.NewAnyWithValue(multiSendMsg)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: "invalid address",
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgExec MsgTransfer": {
			message: func() sdk.Msg {
				msgTransfer := &transfertypes.MsgTransfer{
					Sender:   testAccount1.Address,
					Receiver: testAccount2.Address,
					Token:    coin,
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgTransfer)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: testAccount1.Address,
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgExec MsgTransfer invalid grantee": {
			message: func() sdk.Msg {
				msgTransfer := &transfertypes.MsgTransfer{
					Sender:   testAccount1.Address,
					Receiver: testAccount2.Address,
					Token:    coin,
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgTransfer)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: "invalid address",
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgExec MsgDepositForBurn": {
			message: func() sdk.Msg {
				msgDepositForBurn := &cctptypes.MsgDepositForBurn{
					From:      testAccount2.Address,
					BurnToken: "uusdc",
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgDepositForBurn)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: testAccount1.Address,
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgExec MsgDepositForBurn invalid grantee": {
			message: func() sdk.Msg {
				msgDepositForBurn := &cctptypes.MsgDepositForBurn{
					From:      testAccount2.Address,
					BurnToken: "uusdc",
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgDepositForBurn)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: "invalid address",
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
		"msgExec MsgDepositForBurnWithCaller": {
			message: func() sdk.Msg {
				msgDepositForBurn := &cctptypes.MsgDepositForBurnWithCaller{
					From:      testAccount2.Address,
					BurnToken: "uusdc",
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgDepositForBurn)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: testAccount1.Address,
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: true,
			expectedError:       types.ErrUnauthorized,
		},
		"msgExec MsgDepositForBurnWithCaller invalid grantee": {
			message: func() sdk.Msg {
				msgDepositForBurn := &cctptypes.MsgDepositForBurnWithCaller{
					From:      testAccount2.Address,
					BurnToken: "uusdc",
				}
				msgSendAny, err := codectypes.NewAnyWithValue(msgDepositForBurn)
				require.NoError(t, err)
				msg := &authz.MsgExec{
					Grantee: "invalid address",
					Msgs:    []*codectypes.Any{msgSendAny},
				}
				return msg
			}(),
			blacklistSendAndRec: false,
			testInvalidAddress:  true,
			expectedError:       bech32.ErrInvalidCharacter(32),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// ARRANGE: setup tokenfactory and isBlacklisted decorator
			ftf, ctx := mocks.FiatTokenfactoryKeeper()
			ftf.SetMintingDenom(ctx, types.MintingDenom{Denom: "uusdc"})
			ftf.SetPaused(ctx, types.Paused{Paused: false})
			cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
			ad := fiattokenfactory.NewIsBlacklistedDecorator(ftf)

			// ARRANGE: Build transactions with specific test case message
			builder, err := newMockTxBuilder(cdc)
			require.NoError(t, err)
			if tc.message != nil {
				err = builder.SetMsgs(tc.message)
				require.NoError(t, err)
			}
			tx := builder.GetTx()

			// ACT: Run transaction through ante handler without blacklisting
			_, err = ad.AnteHandle(ctx, tx, true, mockNext)

			// ASSERT: If we are testing for an invalid address, raise error here
			if tc.testInvalidAddress {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			if tc.blacklistSendAndRec {
				// ARRANGE: Blacklist sender account
				ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: testAccount1.AddressBz})

				// ACT: Run transaction through ante handler while sender is blacklisted
				_, err = ad.AnteHandle(ctx, tx, true, mockNext)

				// ASSERT: Assert that the unauthorized error is raised
				require.ErrorIs(t, err, tc.expectedError)

				// ARRANGE: Un-blacklist sender and blacklist receiver
				ftf.RemoveBlacklisted(ctx, testAccount1.AddressBz)
				ftf.SetBlacklisted(ctx, types.Blacklisted{AddressBz: testAccount2.AddressBz})

				// ACT: Run transaction through ante handler while receiver is blacklisted
				_, err = ad.AnteHandle(ctx, tx, true, mockNext)

				// ASSERT: Assert that the unauthorized error is raised
				require.ErrorIs(t, err, tc.expectedError)
			}
		})
	}
}

func mockNext(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
	return ctx, nil
}

func newMockTxBuilder(cdc codec.Codec) (client.TxBuilder, error) {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	builder := txConfig.NewTxBuilder()
	privKey := secp256k1.GenPrivKeyFromSecret([]byte("test"))
	pubKey := privKey.PubKey()
	return builder, builder.SetSignatures(
		signingtypes.SignatureV2{
			PubKey:   pubKey,
			Sequence: 0,
			Data:     &signingtypes.SingleSignatureData{},
		},
	)
}
