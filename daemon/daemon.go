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

type DaemonRpcClient struct {
	connData RpcConnection // Monero daemon rpc connection data
	timeout  time.Duration // Sync timeout
	lbh      uint64        // Last block height synchronized
	dlh      DaemonListenerHandler
	httpcl   http.Client
}

func (c *DaemonRpcClient) sendRequest(method string, path string, body io.Reader) (*http.Response, error) {
	url, err := url.JoinPath(c.connData.host.String(), path)
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
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_block_count
func (c *DaemonRpcClient) GetBlockCount() (*MoneroRpcGenericResponse[GetBlockCountResult], error) {
	req := &MoneroRpcRequest[GetBlockCountParams]{defaultMoneroRpcHeader, "get_block_count", GetBlockCountParams{}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetBlockCountResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// on_get_block_hash
func (c *DaemonRpcClient) OnGetBlockHash(height uint64) (*MoneroRpcGenericResponse[OnGetBlockHashResult], error) {
	req := &MoneroRpcRequest[OnGetBlockHashParams]{defaultMoneroRpcHeader, "on_get_block_hash", [1]uint64{height}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[OnGetBlockHashResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_block_template
func (c *DaemonRpcClient) GetBlockTemplate(wallet string, reverseSize uint64) (*MoneroRpcGenericResponse[GetBlockTemplateResult], error) {
	req := &MoneroRpcRequest[GetBlockTemplateParams]{defaultMoneroRpcHeader, "get_block_template", GetBlockTemplateParams{wallet, reverseSize}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetBlockTemplateResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_last_block_header
func (c *DaemonRpcClient) GetLastBlockHeader(fillPowHash bool) (*MoneroRpcGenericResponse[GetLastBlockHeaderResult], error) {
	req := &MoneroRpcRequest[GetLastBlockHeaderParams]{defaultMoneroRpcHeader, "get_last_block_header", GetLastBlockHeaderParams{fillPowHash}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetLastBlockHeaderResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_block_header_by_hash
func (c *DaemonRpcClient) GetBlockHeaderByHash(fillPowHash bool, hash string) (*MoneroRpcGenericResponse[GetBlockHeaderByHashResult], error) {
	req := &MoneroRpcRequest[GetBlockHeaderByHashParams]{defaultMoneroRpcHeader, "get_block_header_by_hash", GetBlockHeaderByHashParams{GetLastBlockHeaderParams{fillPowHash}, hash}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetBlockHeaderByHashResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_block_header_by_height
func (c *DaemonRpcClient) GetBlockHeaderByHeight(fillPowHash bool, height uint64) (*MoneroRpcGenericResponse[GetBlockHeaderByHeightResult], error) {
	req := &MoneroRpcRequest[GetBlockHeaderByHeightParams]{defaultMoneroRpcHeader, "get_block_header_by_height", GetBlockHeaderByHeightParams{GetLastBlockHeaderParams{fillPowHash}, height}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetBlockHeaderByHeightResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_block_headers_range
func (c *DaemonRpcClient) GetBlockHeadersRange(fillPowHash bool, startHeight uint64, endHeight uint64) (*MoneroRpcGenericResponse[GetBlockHeadersRangeResult], error) {
	req := &MoneroRpcRequest[GetBlockHeadersRangeParams]{defaultMoneroRpcHeader, "get_block_headers_range", GetBlockHeadersRangeParams{GetLastBlockHeaderParams{fillPowHash}, startHeight, endHeight}}

	res, err := getResultFromDaemonRpc[MoneroRpcGenericResponse[GetBlockHeadersRangeResult]](c, http.MethodPost, "/json_rpc", req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

func CreateDaemonRpcClient(c *RpcConnection, timeout time.Duration, dlh DaemonListenerHandler) *DaemonRpcClient {
	return &DaemonRpcClient{*c, timeout, 0, dlh, http.Client{}}
}
