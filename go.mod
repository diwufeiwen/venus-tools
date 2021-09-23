module venus-tools

go 1.16

require (
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-bitfield v0.2.4
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/go-state-types v0.1.1-0.20210506134452-99b279731c48
	github.com/filecoin-project/specs-actors v0.9.14
	github.com/filecoin-project/venus v1.0.5-0.20210907075013-1925e5d6efdd
	github.com/filecoin-project/venus-messager v1.1.1
	github.com/ipfs-force-community/venus-common-utils v0.0.0-20210714054928-2042a9040759
	github.com/ipfs/go-log/v2 v2.3.0
	github.com/urfave/cli/v2 v2.3.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20210219115102-f37d292932f2
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/ipfs/go-ipfs-cmds => github.com/ipfs-force-community/go-ipfs-cmds v0.6.1-0.20210521090123-4587df7fa0ab
