package daemon

import (
	"net/url"
)

type RpcConnection struct {
	host     url.URL
	username string
	password string
}

func NewRpcConnection(host *url.URL, username, password string) *RpcConnection {
	return &RpcConnection{*host, username, password}
}
