Go Monero RPC Client
====================

<p align="center">
<img src="https://github.com/omani/go-monero-rpc-client/raw/master/media/img/monero_gopher.png" alt="Monero Gopher" width="200" />
</p>

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
./monero-wallet-rpc --wallet-file /home/$user/stagenetwallet/stagenetwallet --daemon-address YOUR_STAGENET_NODE:38081 --stagenet --rpc-bind-port 6061 --password 'mystagenetwalletpassword' --disable-rpc-login
```
You can either run your own stagenet server for testing purposes or select a remote stagenet node from eg. https://monero.fail/?chain=monero&network=stagenet (not associated with this site. I found it on google).

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
./monero-wallet-rpc --wallet-file /home/$user/stagenetwallet/stagenetwallet --daemon-address YOUR_STAGENET_NODE:38081 --stagenet --rpc-bind-port 6061 --password 'mystagenetwalletpassword' --rpc-login test:testpass
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

As of now, only the wallet RPC has been implemented.

# Contribution
* You can fork this, extend it and contribute back.
* You can contribute with pull requests.

# Donations
You can make me happy by donating Bitcoin to the following address:
```
bc1qgezvfp4s0xme8pdv6aaqu9ayfgnv4mejdlv3tx
```

# LICENSE
MIT License
