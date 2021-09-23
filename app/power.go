package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"

	typegen "github.com/whyrusleeping/cbor-gen"
	"venus-tools/types"

	"github.com/filecoin-project/venus/pkg/specactors"
	types2 "github.com/filecoin-project/venus/pkg/types"
)

var powerCmd = &cli.Command{
	Name:  "power",
	Usage: "Some tools related to power",
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
		terminateSectors,
	},
}

var terminateSectors = &cli.Command{
	Name:      "terminateSectors",
	Usage:     "terminate sectors for specified miner,eg ./venus-tool --msg-api=/ip4/<msg_ip>/tcp/39812 --token=<msg_token> power --from=<wallet addr> --nonce=0 terminateSectors <MINER_ID> 0 <deadline index> <partition index> [sid_01,sid_02]",
	Action: func(cctx *cli.Context) error {
		deadlineIndex, err := strconv.ParseInt(cctx.Args().Get(2), 10, 64)
		if err != nil {
			panic(err)
		}
		partitionIndex, err := strconv.ParseInt(cctx.Args().Get(3), 10, 64)
		if err != nil {
			panic(err)
		}

		sectorIDs := strings.Split(cctx.Args().Get(4), ",")
		bf := bitfield.New()

		for _, sidStr := range sectorIDs {
			sid, err := strconv.ParseUint(sidStr, 10, 64)
			if err != nil {
				return err
			}
			bf.Set(sid)
		}

		params := &miner.TerminateSectorsParams{Terminations: []miner.TerminationDeclaration{miner.TerminationDeclaration{
			Deadline:  uint64(deadlineIndex),
			Partition: uint64(partitionIndex),
			Sectors:   bf,
		}}}

		paramBytes, err := serializeParamsAndSend(params)
		if err != nil {
			return err
		}

		return send(cctx, paramBytes, builtin.MethodsMiner.TerminateSectors)
	},
}

func serializeParamsAndSend(params typegen.CBORMarshaler) ([]byte, error) {
	paramBytes, err := specactors.SerializeParams(params)
	if err != nil {
		return nil, err
	}
	return paramBytes, nil
}

func send(cctx *cli.Context, params []byte, methodNum abi.MethodNum) error {
	from := cctx.String("from")
	if from == "" {
		return fmt.Errorf("from is empty")
	}
	fromAddr, err := address.NewFromString(from)
	if err != nil {
		return err
	}

	toAddr, err := address.NewFromString(cctx.Args().Get(0))
	if err != nil {
		return err
	}

	value, err := types.ParseFIL(cctx.Args().Get(1))
	if err != nil {
		return err
	}

	api, closer, err := NewMessageRPC(&msgConfig)
	if err != nil {
		return err
	}
	defer closer()

	nonce := cctx.Uint64("nonce")

	uid, err := api.PushMessage(cctx.Context, &types2.UnsignedMessage{
		Version: 0,
		To:      toAddr,
		From:    fromAddr,
		Nonce:   nonce,
		Value:   types.BigInt(value),
		Method:  methodNum,
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
	fmt.Println("success!")

	return nil
}
