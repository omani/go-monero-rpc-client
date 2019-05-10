# go-monero-rpc-client

A client implementation for the Monero wallet and daemon RPC written in go.
This package is inspired by https://github.com/gabstv/go-monero.

## Wallet RPC Client

[![GoDoc](https://godoc.org/github.com/omani/go-monero-rpc-client/wallet?status.svg)](https://godoc.org/github.com/omani/go-monero-rpc-client/wallet)

### Monero RPC Version
The ```go-monero-rpc-client/wallet``` package is the RPC client for version `v1.3` of the [Monero Wallet RPC](https://www.getmonero.org/resources/developer-guides/wallet-rpc.html).

### Installation

```sh
go get -u github.com/omani/go-monero-rpc-client
```

#### Spawn the monero-wallet-rpc daemon (without rpc login):

```sh
./monero-wallet-rpc --wallet-file /home/$user/stagenetwallet/stagenetwallet --daemon-address pool.cloudissh.com:38081 --stagenet --rpc-bind-port 6061 --password 'mystagenetwalletpassword' --disable-rpc-login
```
You can use our remote node for the stagenet running at pool.cloudissh.com port `38081`.

#### Go code:

```Go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/omani/go-monero-rpc-client/wallet"
)

func checkerr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	// Start a wallet client instance
	client := wallet.New(wallet.Config{
      Address: "http://127.0.0.1:6061/json_rpc",
	})

	// check wallet balance
	resp, err := client.GetBalance(&wallet.RequestGetBalance{AccountIndex: 0})
	checkerr(err)
	res, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Print(string(res))

	// get incoming transfers
	resp1, err := client.GetTransfers(&wallet.RequestGetTransfers{
		AccountIndex: 0,
		In:           true,
	})
	checkerr(err)
	for _, in := range resp1.In {
		res, _ := json.MarshalIndent(in, "", "\t")
		fmt.Print(string(res))
	}
}
```

### Spawn the monero-wallet-rpc daemon (with rpc login):

```sh
./monero-wallet-rpc --wallet-file /home/$user/stagenetwallet/stagenetwallet --daemon-address pool.cloudissh.com:38081 --stagenet --rpc-bind-port 6061 --password 'mystagenetwalletpassword' --rpc-login test:testpass
```

#### Go code:

```Go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/omani/go-monero-rpc-client/wallet"
)

func checkerr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
  t := httpdigest.New("test", "testpass")

	// Start a wallet client instance
	client := wallet.New(wallet.Config{
      Address: "http://127.0.0.1:6061/json_rpc",
      Transport: t,
	})

	// check wallet balance
	resp, err := client.GetBalance(&wallet.RequestGetBalance{AccountIndex: 0})
	checkerr(err)
	res, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Print(string(res))

	// get incoming transfers
	resp1, err := client.GetTransfers(&wallet.RequestGetTransfers{
		AccountIndex: 0,
		In:           true,
	})
	checkerr(err)
	for _, in := range resp1.In {
		res, _ := json.MarshalIndent(in, "", "\t")
		fmt.Print(string(res))
	}
}
```

# Daemon RPC Client

As of now, only the wallet RPC has been implemented. The daemon RPC will follow very soon.

# Contribution
* You can fork this, extend it and contribute back.
* You can contribute with pull requests.

# Donations
I love Monero (XMR) and building applications for and on top of Monero.

You can make me happy by donating Monero to the following address:

```
89woiq9b5byQ89SsUL4Bd66MNfReBrTwNEDk9GoacgESjfiGnLSZjTD5x7CcUZba4PBbE3gUJRQyLWD4Akz8554DR4Lcyoj
```

# LICENSE
MIT License