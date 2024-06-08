package main

import (
	"log"
	"net/url"
	"time"

	"github.com/monero-ecosystem/go-monero-rpc-client/daemon"
)

func main() {
	host, err := url.Parse("http://node.sethforprivacy.com:18089")
	if err != nil {
		log.Printf(err.Error())
		return
	}

	d := daemon.CreateDaemonRpcClient(daemon.NewRpcConnection(*host, "", ""), 5*time.Second, &daemon.DaemonListenerHandlerAsync{})
	_, err = d.GetBlockHeadersRange(false, 1545999, 1546000)
	if err != nil {
		log.Println(err.Error())
	}
}
