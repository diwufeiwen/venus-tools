package main

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"

	"github.com/ipfs-force-community/venus-common-utils/apiinfo"

	api "github.com/filecoin-project/venus/venus-shared/api/messager"

	fullNodeApi "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

type VenusConfig struct {
	NodeUrl    string
	MessageUrl string
	Token      string
}

func NewMessageRPC(cfg *VenusConfig) (api.IMessager, jsonrpc.ClientCloser, error) {
	apiInfo := apiinfo.APIInfo{
		Addr:  cfg.MessageUrl,
		Token: []byte(cfg.Token),
	}

	addr, err := apiInfo.DialArgs("v0")
	if err != nil {
		return nil, nil, err
	}

	return api.NewIMessagerRPC(context.Background(), addr, apiInfo.AuthHeader())
}

func NewFullNodeRPC(cfg *VenusConfig) (fullNodeApi.FullNode, jsonrpc.ClientCloser, error) {
	apiInfo := apiinfo.APIInfo{
		Addr:  cfg.MessageUrl,
		Token: []byte(cfg.Token),
	}

	addr, err := apiInfo.DialArgs("v0")
	if err != nil {
		return nil, nil, err
	}
	return fullNodeApi.NewFullNodeRPC(context.Background(), addr, apiInfo.AuthHeader())
}
