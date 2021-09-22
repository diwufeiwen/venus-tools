package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-state-types/big"
	"strconv"
	"strings"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"

	"github.com/ipfs-force-community/venus-common-utils/apiinfo"

	typegen "github.com/whyrusleeping/cbor-gen"
	"venus-tools/types"

	"github.com/filecoin-project/venus/pkg/specactors"
	types2 "github.com/filecoin-project/venus/pkg/types"

	"github.com/filecoin-project/venus-messager/api/client"
	types3 "github.com/filecoin-project/venus-messager/types"
)

type TaskType string

const (
	TerminateSectors TaskType = "terminateSectors"
)

type IMessager interface {
	WalletHas(ctx context.Context, addr address.Address) (bool, error)
	WaitMessage(ctx context.Context, id string, confidence uint64) (*types3.Message, error)
	PushMessage(ctx context.Context, msg *types2.UnsignedMessage, meta *types3.MsgMeta) (string, error)
	PushMessageWithId(ctx context.Context, id string, msg *types2.UnsignedMessage, meta *types3.MsgMeta) (string, error)
	GetMessageByUid(ctx context.Context, id string) (*types3.Message, error)
}

func NewMessageRPC(cfg *MessagerConfig) (IMessager, jsonrpc.ClientCloser, error) {
	apiInfo := apiinfo.APIInfo{
		Addr:  cfg.Url,
		Token: []byte(cfg.Token),
	}

	addr, err := apiInfo.DialArgs("v0")
	if err != nil {
		return nil, nil, err
	}

	return client.NewMessageRPC(context.Background(), addr, apiInfo.AuthHeader())
}

func serializeParamsAndSend(params typegen.CBORMarshaler) ([]byte, error) {
	paramBytes, err := specactors.SerializeParams(params)
	if err != nil {
		return nil, err
	}
	return paramBytes, nil
}

func TerminateSectorsProc(from, minerAddr, value, di, pi, sids string, nonce uint64) (string, error) {
	var fromAddr address.Address
	if from == "" {
		return "", fmt.Errorf("from is empty")
	} else {
		addr, err := address.NewFromString(from)
		if err != nil {
			return "", err
		}

		fromAddr = addr
	}

	toAddr, err := address.NewFromString(minerAddr)
	if err != nil {
		return "", err
	}

	val, err := types.ParseFIL(value)
	if err != nil {
		return "", err
	}

	deadlineIndex, err := strconv.ParseInt(di, 10, 64)
	if err != nil {
		panic(err)
	}
	partitionIndex, err := strconv.ParseInt(pi, 10, 64)
	if err != nil {
		panic(err)
	}

	sectorIDs := strings.Split(sids, ",")
	bf := bitfield.New()

	for k := range sectorIDs {
		sid, err := strconv.ParseUint(sectorIDs[k], 10, 64)
		if err != nil {
			return "", err
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
		return "", err
	}

	return send(fromAddr, toAddr, val, paramBytes, abi.MethodNum(9), nonce)
}

func send(from, to address.Address, val types.FIL, params []byte, methodNum abi.MethodNum, nonce uint64) (string, error) {
	api, closer, err := NewMessageRPC(&msgConfig)
	if err != nil {
		return "", err
	}
	defer closer()

	uid, err := api.PushMessage(context.Background(), &types2.UnsignedMessage{
		Version: 0,
		To:      to,
		From:    from,
		Nonce:   nonce,
		Value:   types.BigInt(val),
		Method:  methodNum,
		Params:  params,
	}, &types3.MsgMeta{
		ExpireEpoch:       0,
		GasOverEstimation: 0,
		//todo give a maxFee or for nil
		MaxFee: big.Zero(),
	})

	return uid, nil
}
