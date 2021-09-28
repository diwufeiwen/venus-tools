package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/builtin"

	"venus-tools/types"

	types2 "github.com/filecoin-project/venus/pkg/types"
)

var filCmd = &cli.Command{
	Name:  "fil",
	Usage: "Some tools related to fil",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Required: true,
			Usage:    "specify the account to send message from",
		},
		&cli.Uint64Flag{
			Name:  "nonce",
			Usage: "self-specified value to fill in missing nonce",
			Value: 0,
		},
	},
	Subcommands: []*cli.Command{
		sendCmd,
	},
}

var sendCmd = &cli.Command{
	Name:      "send",
	Usage:     "Send funds between accounts",
	ArgsUsage: "[targetAddress] [amount]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() != 2 {
			return ShowHelp(cctx, fmt.Errorf("'send' expects two arguments, target and amount"))
		}

		toAddr, err := address.NewFromString(cctx.Args().Get(0))
		if err != nil {
			return err
		}

		value, err := types.ParseFIL(cctx.Args().Get(1))
		if err != nil {
			return err
		}

		from := cctx.String("from")
		if from == "" {
			return fmt.Errorf("from is empty")
		}
		fromAddr, err := address.NewFromString(from)
		if err != nil {
			return err
		}

		nonce := cctx.Uint64("nonce")
		msg := &types2.UnsignedMessage{
			From:  fromAddr,
			To:    toAddr,
			Value: types.BigInt(value),

			Method: builtin.MethodSend,
			Nonce:   nonce,
		}

		api, closer, err := NewMessageRPC(&msgConfig)
		if err != nil {
			return err
		}
		defer closer()

		uid, err := api.PushMessage(cctx.Context, msg, nil)

		if err != nil {
			return err
		}
		fmt.Printf("msg uid is : %s, search for the processing result from venus-messager ...\n", uid)

		return  nil
	},
}
