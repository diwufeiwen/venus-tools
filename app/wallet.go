package main

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/venus/venus-shared/actors"
	"github.com/filecoin-project/venus/venus-shared/actors/builtin/market"
	types2 "github.com/filecoin-project/venus/venus-shared/types"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"venus-tools/types"
)

var marketCmds = &cli.Command{
	Name:  "market",
	Usage: "Interact with market balances",
	Subcommands: []*cli.Command{
		walletMarketAdd,
	},
}

var walletMarketAdd = &cli.Command{
	Name:      "add",
	Usage:     "Add funds to the Storage Market Actor",
	ArgsUsage: "<amount>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "from",
			Usage:   "Specify address to move funds from, otherwise it will use the default wallet address",
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Market address to move funds to (account or miner actor address, defaults to --from address)",
			Aliases: []string{"a"},
		},
	},
	Action: func(cctx *cli.Context) error {
		api, closer, err := NewMessageRPC(&msgConfig)
		if err != nil {
			return err
		}
		defer closer()

		afmt := NewAppFmt(cctx.App)

		// Get amount param
		if !cctx.Args().Present() {
			return fmt.Errorf("must pass amount to add")
		}
		f, err := types.ParseFIL(cctx.Args().First())
		if err != nil {
			return xerrors.Errorf("parsing 'amount' argument: %w", err)
		}

		amt := abi.TokenAmount(f)


		// Get address param
		from := cctx.String("from")
		if from == "" {
			return fmt.Errorf("from is empty")
		}
		fromAddr, err := address.NewFromString(from)
		if err != nil {
			return err
		}

		addrStr := cctx.String("address")
		if addrStr  == "" {
			return fmt.Errorf("from is empty")
		}
		addr, err := address.NewFromString(addrStr)
		if err != nil {
			return err
		}

		// Add balance to market actor
		fmt.Printf("Submitting Add Balance message for amount %s for address %s\n", types.FIL(amt), fromAddr)
		params, err := actors.SerializeParams(&addr)
		if err != nil {
			return  err
		}
		uid, err := api.PushMessage(cctx.Context, &types2.Message{
			Version: 0,
			To:      market.Address,
			From:    fromAddr,
			Nonce:   0,
			Value:   amt,
			Method:  market.Methods.AddBalance,
			Params:  params,
		}, nil)

		if err != nil {
			return err
		}
		fmt.Printf("msg uid is : %s, waiting for the processing result ...\n", uid)

		mw, err := api.WaitMessage(cctx.Context, uid, MessageConfidence)
		if err != nil {
			return fmt.Errorf("waiting for worker init: %w", err)
		}

		if mw.Receipt.ExitCode != 0 {
			return  fmt.Errorf("msg run failed, exit code %d", mw.Receipt.ExitCode)
		}
		afmt.Printf("AddBalance message cid: %s\n", mw.Cid())

		return nil
	},
}
