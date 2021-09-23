package main

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"

	"github.com/ipfs-force-community/venus-common-utils/apiinfo"

	types2 "github.com/filecoin-project/venus/pkg/types"

	"github.com/filecoin-project/venus-messager/api/client"
	types3 "github.com/filecoin-project/venus-messager/types"
)

type MessagerConfig struct {
	Url   string
	Token string
}

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
