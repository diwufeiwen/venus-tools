package main

import (
	"flag"
	"fmt"
)

type MessagerConfig struct {
	Url   string
	Token string
}

var msgConfig MessagerConfig

func main() {
	var (
		taskType, from   string
		minerAddr, value string
		nonce            uint64
		di, pi, sids     string

		url, token string
	)

	flag.StringVar(&taskType, "task", "terminateSectors", "task type: terminateSectors")

	// TerminateSectors
	flag.StringVar(&from, "from", "", "address to send message from")
	flag.Uint64Var(&nonce, "nonce", 0, "specify the nonce to use")
	flag.StringVar(&minerAddr, "minerAddr", "", "miner id")
	flag.StringVar(&value, "value", "", "value to send with message in FIL")
	flag.StringVar(&di, "di", "", "deadline Index")
	flag.StringVar(&pi, "pi", "", "partition Index")
	flag.StringVar(&sids, "sids", "", "sectors split by comma")

	flag.StringVar(&url, "url", "", "venus-messager api, eg. /ip4/<ip>/tcp/<port>")
	flag.StringVar(&token, "token", "", "venus-messager token")

	flag.Parse()

	msgConfig = MessagerConfig{Url: url, Token: token}

	switch TaskType(taskType) {
	case TerminateSectors:
		uid, err := TerminateSectorsProc(from, minerAddr, value, di, pi, sids, nonce)
		if err == nil {
			fmt.Printf("msg cid is %s\n", uid)
		} else {
			fmt.Printf("proc [%s] err: %s\n", taskType, err)

		}
	default:
		fmt.Printf("Invalid task type: %s\n", taskType)
	}
}
