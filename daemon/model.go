package daemon

type MoneroRpcRequestParams interface {
	OnGetBlockHashParams | EmptyMoneroRpcParams
}

type MoneroRpcResponse interface {
	MoneroRpcGenericResponse[GetBlockCountResult] | MoneroRpcGenericResponse[OnGetBlockHashResult] | GetHeightResponse
}

type MoneroRpcResponseResult interface {
	GetBlockCountResult | OnGetBlockHashResult
}

type JsonRpcHeader struct {
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}
type JsonRpcFooter struct {
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

type MoneroRpcRequest[T MoneroRpcRequestParams] struct {
	JsonRpcHeader
	Method string `json:"method"`
	Params T      `json:"params"`
}
type EmptyMoneroRpcParams string

type MoneroRpcGenericResponse[T MoneroRpcResponseResult] struct {
	JsonRpcHeader
	Result T `json:"result"`
}

/**
	BASIC RPC METHODS
**/

// get_block_count
type GetBlockCountResult struct {
	Count uint64 `json:"count"`
	JsonRpcFooter
}

// on_get_block_hash
type OnGetBlockHashParams [1]int

type OnGetBlockHashResult string

/**
	OTHER RPC METHODS
**/

// get_height
type GetHeightResponse struct {
	Hash   string `json:"hash"`
	Height uint64 `json:"height"`
	JsonRpcFooter
}
