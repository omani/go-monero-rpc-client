package daemon

type RpcConnection struct {
	host     string
	username string
	password string
}

func NewRpcConnection(host, username, password string) *RpcConnection {
	return &RpcConnection{host, username, password}
}
