package cli

import (
	"strconv"

	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdUpdateOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-owner [address]",
		Short: "Broadcast message update-owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateOwner{
				From:    clientCtx.GetFromAddress().String(),
				Address: argAddress,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
