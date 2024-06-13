package daemon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/icholy/digest"
	"github.com/monero-ecosystem/go-monero-rpc-client/util"
)

var defaultMoneroRpcHeader = JsonRpcHeader{Id: "0", Jsonrpc: "2.0"}

type IDaemonRpcClient interface {
	SetRpcConnection(connection *RpcConnection)

	/**
		JSON RPC METHODS
	**/

	// get_block_count
	GetBlockCount() (*JsonRpcGenericResponse[GetBlockCountResult], error)
	// on_get_block_hash
	OnGetBlockHash(height uint64) (*JsonRpcGenericResponse[OnGetBlockHashResult], error)
	// get_block_template
	GetBlockTemplate(wallet string, reverseSize uint64) (*JsonRpcGenericResponse[GetBlockTemplateResult], error)
	// get_last_block_header
	GetLastBlockHeader(fillPowHash bool) (*JsonRpcGenericResponse[GetBlockHeaderResult], error)
	// get_block_header_by_hash
	GetBlockHeaderByHash(fillPowHash bool, hash string) (*JsonRpcGenericResponse[GetBlockHeaderResult], error)
	// get_block_header_by_height
	GetBlockHeaderByHeight(fillPowHash bool, height uint64) (*JsonRpcGenericResponse[GetBlockHeaderResult], error)
	// get_block_headers_range
	GetBlockHeadersRange(fillPowHash bool, startHeight uint64, endHeight uint64) (*JsonRpcGenericResponse[GetBlockHeadersRangeResult], error)
	// get_block
	GetBlockByHeight(fillPowHash bool, height uint64) (*JsonRpcGenericResponse[GetBlockResult], error)
	// get_block
	GetBlockByHash(fillPowHash bool, hash string) (*JsonRpcGenericResponse[GetBlockResult], error)
	// get_fee_estimate
	GetFeeEstimate() (*JsonRpcGenericResponse[GetFeeEstimateResult], error)
	// get_version
	GetVersion() (*JsonRpcGenericResponse[GetVersionResult], error)
	// get_info
	GetInfo() (*JsonRpcGenericResponse[GetInfoResult], error)

	/**
		OTHER RPC METHODS
	**/

	// get_height
	GetCurrentHeight() (*GetHeightResponse, error)
	// get_transaction_pool
	GetTransactionPool() (*GetTransactionPoolResponse, error)
	// get_transactions
	GetTransactions(txHashes []string, decodeAsJson bool, prune bool, split bool) (*GetTransactionsResponse, error)
}

type DaemonRpcClient struct {
	connData RpcConnection
	httpcl   *http.Client
}

func (c *DaemonRpcClient) SetRpcConnection(connection *RpcConnection) {
	c.connData = *connection
	c.httpcl.Transport = &digest.Transport{
		Username: connection.username,
		Password: connection.password,
	}
}

func (c *DaemonRpcClient) sendRequest(method string, path string, body io.Reader) (*http.Response, error) {
	url := c.connData.host.Scheme + "://" + c.connData.host.Host + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	return c.httpcl.Do(req)
}

func getResultFromDaemonRpc[R MoneroRpcResponse, B MoneroRpcRequestBody](c *DaemonRpcClient, req *MoneroRpcRequest[B]) (*R, error) {
	var data []byte
	var err error

	body := req.Body
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	res, err := c.sendRequest(http.MethodPost, req.Endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	result, err := util.ParseResponse[R](res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewDaemonRpcClient(connection *RpcConnection) IDaemonRpcClient {
	return &DaemonRpcClient{
		connData: *connection,
		httpcl: &http.Client{
			Transport: &digest.Transport{
				Username: connection.username,
				Password: connection.password,
			},
		},
	}
}

/**
	JSON RPC METHODS
**/

// get_block_count
func (c *DaemonRpcClient) GetBlockCount() (*JsonRpcGenericResponse[GetBlockCountResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockCountParams]{defaultMoneroRpcHeader, "get_block_count", GetBlockCountParams{}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockCountParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockCountResult]](c, req)
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
func (c *DaemonRpcClient) OnGetBlockHash(height uint64) (*JsonRpcGenericResponse[OnGetBlockHashResult], error) {
	reqBody := &JsonRpcGenericRequestBody[OnGetBlockHashParams]{defaultMoneroRpcHeader, "on_get_block_hash", [1]uint64{height}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[OnGetBlockHashParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[OnGetBlockHashResult]](c, req)
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
func (c *DaemonRpcClient) GetBlockTemplate(wallet string, reverseSize uint64) (*JsonRpcGenericResponse[GetBlockTemplateResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockTemplateParams]{defaultMoneroRpcHeader, "get_block_template", GetBlockTemplateParams{wallet, reverseSize}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockTemplateParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockTemplateResult]](c, req)
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
func (c *DaemonRpcClient) GetLastBlockHeader(fillPowHash bool) (*JsonRpcGenericResponse[GetBlockHeaderResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockHeaderDefaultParams]{defaultMoneroRpcHeader, "get_last_block_header", GetBlockHeaderDefaultParams{fillPowHash}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockHeaderDefaultParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockHeaderResult]](c, req)
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
func (c *DaemonRpcClient) GetBlockHeaderByHash(fillPowHash bool, hash string) (*JsonRpcGenericResponse[GetBlockHeaderResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockHeaderByHashParams]{defaultMoneroRpcHeader, "get_block_header_by_hash", GetBlockHeaderByHashParams{GetBlockHeaderDefaultParams{fillPowHash}, hash}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockHeaderByHashParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockHeaderResult]](c, req)
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
func (c *DaemonRpcClient) GetBlockHeaderByHeight(fillPowHash bool, height uint64) (*JsonRpcGenericResponse[GetBlockHeaderResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockHeaderByHeightParams]{defaultMoneroRpcHeader, "get_block_header_by_height", GetBlockHeaderByHeightParams{GetBlockHeaderDefaultParams{fillPowHash}, height}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockHeaderByHeightParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockHeaderResult]](c, req)
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
func (c *DaemonRpcClient) GetBlockHeadersRange(fillPowHash bool, startHeight uint64, endHeight uint64) (*JsonRpcGenericResponse[GetBlockHeadersRangeResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockHeadersRangeParams]{defaultMoneroRpcHeader, "get_block_headers_range", GetBlockHeadersRangeParams{GetBlockHeaderDefaultParams{fillPowHash}, startHeight, endHeight}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockHeadersRangeParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockHeadersRangeResult]](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

func fillMissedFieldHelper[T any](json *string, field *T) error {
	val, err := util.ParseJsonString[T](json)
	if err != nil {
		return err
	}

	*field = *val

	return nil
}

// get_block
func fillBlockDetailsHelper(res *JsonRpcGenericResponse[GetBlockResult]) (*JsonRpcGenericResponse[GetBlockResult], error) {
	if err := fillMissedFieldHelper(&res.Result.Json, &res.Result.BlockDetails); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return res, nil
}

// get_block
func (c *DaemonRpcClient) GetBlockByHeight(fillPowHash bool, height uint64) (*JsonRpcGenericResponse[GetBlockResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockByHeightParams]{defaultMoneroRpcHeader, "get_block", GetBlockByHeightParams{GetBlockHeaderDefaultParams{fillPowHash}, height}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockByHeightParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockResult]](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return fillBlockDetailsHelper(res)
}

// get_block
func (c *DaemonRpcClient) GetBlockByHash(fillPowHash bool, hash string) (*JsonRpcGenericResponse[GetBlockResult], error) {
	reqBody := &JsonRpcGenericRequestBody[GetBlockByHashParams]{defaultMoneroRpcHeader, "get_block", GetBlockByHashParams{GetBlockHeaderDefaultParams{fillPowHash}, hash}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[GetBlockByHashParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetBlockResult]](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return fillBlockDetailsHelper(res)
}

// get_fee_estimate
func (c *DaemonRpcClient) GetFeeEstimate() (*JsonRpcGenericResponse[GetFeeEstimateResult], error) {
	reqBody := &JsonRpcGenericRequestBody[EmptyMoneroRpcParams]{defaultMoneroRpcHeader, "get_fee_estimate", EmptyMoneroRpcParams{}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[EmptyMoneroRpcParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetFeeEstimateResult]](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_version
func (c *DaemonRpcClient) GetVersion() (*JsonRpcGenericResponse[GetVersionResult], error) {
	reqBody := &JsonRpcGenericRequestBody[EmptyMoneroRpcParams]{defaultMoneroRpcHeader, "get_version", EmptyMoneroRpcParams{}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[EmptyMoneroRpcParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetVersionResult]](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_info
func (c *DaemonRpcClient) GetInfo() (*JsonRpcGenericResponse[GetInfoResult], error) {
	reqBody := &JsonRpcGenericRequestBody[EmptyMoneroRpcParams]{defaultMoneroRpcHeader, "get_info", EmptyMoneroRpcParams{}}
	req := &MoneroRpcRequest[JsonRpcGenericRequestBody[EmptyMoneroRpcParams]]{DEFAULT_MONERO_RPC_ENDPOINT, reqBody}

	res, err := getResultFromDaemonRpc[JsonRpcGenericResponse[GetInfoResult]](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

/**
	OTHER RPC METHODS
**/

// get_height
func (c *DaemonRpcClient) GetCurrentHeight() (*GetHeightResponse, error) {
	req := &MoneroRpcRequest[EmptyMoneroRpcParams]{"/get_height", nil}

	res, err := getResultFromDaemonRpc[GetHeightResponse](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if res.Error.Code != 0 {
		return nil, &res.Error
	}

	return res, nil
}

// get_transaction_pool
func fillGetTransactionPoolHelper(res *GetTransactionPoolResponse) (*GetTransactionPoolResponse, error) {
	for i := range res.Transactions {
		if err := fillMissedFieldHelper(&res.Transactions[i].TxJson, &res.Transactions[i].TxInfo); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	return res, nil
}

// get_transaction_pool
func (c *DaemonRpcClient) GetTransactionPool() (*GetTransactionPoolResponse, error) {
	req := &MoneroRpcRequest[EmptyMoneroRpcParams]{"/get_transaction_pool", nil}

	res, err := getResultFromDaemonRpc[GetTransactionPoolResponse](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return fillGetTransactionPoolHelper(res)
}

func fillGetTransactionsHelper(res *GetTransactionsResponse) (*GetTransactionsResponse, error) {
	for i := range res.Txs {
		if err := fillMissedFieldHelper(&res.Txs[i].AsJson, &res.Txs[i].TxInfo); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	return res, nil
}

// get_transactions
func (c *DaemonRpcClient) GetTransactions(txHashes []string, decodeAsJson bool, prune bool, split bool) (*GetTransactionsResponse, error) {
	reqBody := &GetTransactionsParams{
		TxHashes:     txHashes,
		DecodeAsJson: decodeAsJson,
		Prune:        prune,
		Split:        split,
	}
	req := &MoneroRpcRequest[GetTransactionsParams]{"/get_transactions", reqBody}

	res, err := getResultFromDaemonRpc[GetTransactionsResponse](c, req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if decodeAsJson {
		return fillGetTransactionsHelper(res)
	}

	return res, nil
}
