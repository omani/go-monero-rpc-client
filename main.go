package main

import (
	"time"

	"github.com/monero-ecosystem/go-monero-rpc-client/daemon"
)

func main() {
	d := daemon.Create(daemon.NewRpcConnection("https://node.sethforprivacy.com/", "", ""), 5*time.Second, &daemon.DaemonListenerHandlerAsync{})
	d.GetBlockCount()
}
