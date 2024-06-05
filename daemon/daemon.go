package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/monero-ecosystem/go-monero-rpc-client/util"
)

var defaultMoneroRpcHeader = JsonRpcHeader{"0", "2.0"}

type IDaemonRpcClient interface {
	Connect() error
	Reconnect() error

	StartSync() error
	StopSync()

	SetRpcConnection(c *RpcConnection)
	SetTimeout(timeout uint32)

	sendRequest()

	GetCurrentHeight()
}

type DaemonRpcClient struct {
	// Monero daemon rpc connection data
	connData RpcConnection
	// Sync timeout
	timeout time.Duration
	// Last block height synchronized
	lbh    uint64
	dlh    DaemonListenerHandler
	httpcl http.Client
}

func (c *DaemonRpcClient) sendRequest(method string, path string, body io.Reader) (*http.Response, error) {
	url, err := url.JoinPath(c.connData.host, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	return c.httpcl.Do(req)
}

func getResultFromDaemonRpc[R MoneroRpcResponse, P MoneroRpcRequestParams](c *DaemonRpcClient, method string, path string, body *MoneroRpcRequest[P]) (*R, error) {
	var data []byte
	var err error

	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	res, err := c.sendRequest(method, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	result, err := util.ParseResponse[R](res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *DaemonRpcClient) GetCurrentHeight() (*GetHeightResponse, error) {
	res, err := getResultFromDaemonRpc[GetHeightResponse, EmptyMoneroRpcParams](c, http.MethodGet, "/get_height", nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return res, nil
}

func (c *DaemonRpcClient) GetBlockCount() (*MoneroRpcGenericResponse[GetBlockCountResult], error) {
	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetBlockCountResult]](c, http.MethodPost, "/json_rpc", &MoneroRpcRequest[EmptyMoneroRpcParams]{defaultMoneroRpcHeader, "get_block_count", ""})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return res, nil
}

func Create(c *RpcConnection, timeout time.Duration, dlh DaemonListenerHandler) *DaemonRpcClient {
	return &DaemonRpcClient{*c, timeout, 0, dlh, http.Client{}}
}
