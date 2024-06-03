package blockibc_test

import (
	"testing"

	"github.com/circlefin/noble-fiattokenfactory/utils"
	"github.com/circlefin/noble-fiattokenfactory/utils/mocks"
	fiattokenfactorytypes "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/x/auth/codec"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestBlockIBC(t *testing.T) {

	// ARRANGE: Mock sender and receiver.
	sender, receiver := utils.TestAccount(), utils.TestAccount()
	receiverAddress, _ := codec.NewBech32Codec("osmo").BytesToString(receiver.AddressBz)

	// ARRANGE: Organize table driven test cases.
	testCases := map[string]struct {
		toBlacklist         *utils.Account
		setPaused           bool
		packet              channeltypes.Packet
		expectSuccessfulAck bool
	}{
		"happy path": {
			toBlacklist:         nil,
			setPaused:           false,
			packet:              mockPacket(sender.Address, receiverAddress),
			expectSuccessfulAck: true,
		},
		"malformed ICS-20 packet data": {
			toBlacklist: nil,
			setPaused:   false,
			packet: func() channeltypes.Packet {
				packet := mockPacket(sender.Address, receiverAddress)
				packet.Data = []byte("malformed packet data")
				return packet
			}(),
			expectSuccessfulAck: false,
		},
		"uncontrolled denom": {
			toBlacklist: nil,
			setPaused:   false,
			packet: func() channeltypes.Packet {
				packet := mockPacket(sender.Address, receiverAddress)
				// transfer `ustake` instead of `usdc`
				packet.Data = transfertypes.NewFungibleTokenPacketData(
					"ustake", "1000000", sender.Address, receiverAddress, "",
				).GetBytes()
				return packet
			}(),
			expectSuccessfulAck: true,
		},
		"tokenfactory paused": {
			toBlacklist:         nil,
			setPaused:           true,
			packet:              mockPacket(sender.Address, receiverAddress),
			expectSuccessfulAck: false,
		},
		"blacklisted sender": {
			toBlacklist:         &sender,
			setPaused:           false,
			packet:              mockPacket(sender.Address, receiverAddress),
			expectSuccessfulAck: false,
		},
		"blacklisted receiver": {
			toBlacklist:         &receiver,
			setPaused:           false,
			packet:              mockPacket(sender.Address, receiverAddress),
			expectSuccessfulAck: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			// ARRANGE: Mock middleware stack.
			middleware, ftf, ctx := mocks.BlockIBC()

			// ACT: Set paused and blacklisted state based on test case.
			if tc.setPaused {
				ftf.SetPaused(ctx, fiattokenfactorytypes.Paused{Paused: true})
			}
			if tc.toBlacklist != nil {
				ftf.SetBlacklisted(ctx, fiattokenfactorytypes.Blacklisted{
					AddressBz: tc.toBlacklist.AddressBz,
				})
			}

			// ACT: Receive transfer packet in middleware.
			ack := middleware.OnRecvPacket(ctx, tc.packet, nil)

			// ASSERT: Assert the acknowledgment's success based on the test case.
			var assertBool require.BoolAssertionFunc
			if tc.expectSuccessfulAck {
				assertBool = require.True
			} else {
				assertBool = require.False
			}
			assertBool(t, ack.Success())
		})
	}
}

func mockPacket(sender, receiver string) channeltypes.Packet {
	return channeltypes.NewPacket(
		transfertypes.NewFungibleTokenPacketData(
			"uusdc", "1000000", sender, receiver, "",
		).GetBytes(),
		1,
		transfertypes.PortID,
		"channel-0",
		transfertypes.PortID,
		"channel-0",
		clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: 0,
		},
		1234,
	)
}
