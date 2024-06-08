package daemon

import "encoding/json"

type MoneroRpcRequestParams interface {
	OnGetBlockHashParams |
		GetBlockCountParams |
		GetLastBlockHeaderParams |
		GetBlockTemplateParams |
		EmptyMoneroRpcParams |
		GetBlockHeaderByHashParams |
		GetBlockHeaderByHeightParams |
		GetBlockHeadersRangeParams
}

type MoneroRpcResponse interface {
	MoneroRpcGenericResponse[GetBlockCountResult] |
		MoneroRpcGenericResponse[OnGetBlockHashResult] |
		MoneroRpcGenericResponse[GetBlockTemplateResult] |
		MoneroRpcGenericResponse[GetLastBlockHeaderResult] |
		MoneroRpcGenericResponse[GetBlockHeaderByHashResult] |
		MoneroRpcGenericResponse[GetBlockHeaderByHeightResult] |
		MoneroRpcGenericResponse[GetBlockHeadersRangeResult] |
		GetHeightResponse
}

type MoneroRpcResponseResult interface {
	GetBlockCountResult |
		OnGetBlockHashResult |
		GetBlockTemplateResult |
		GetLastBlockHeaderResult |
		GetBlockHeaderByHashResult |
		GetBlockHeaderByHeightResult |
		GetBlockHeadersRangeResult
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
type EmptyMoneroRpcParams struct{}

type MoneroRpcGenericResponse[T MoneroRpcResponseResult] struct {
	JsonRpcHeader
	Result T              `json:"result"`
	Error  MoneroRpcError `json:"error"`
}

type MoneroRpcError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (e *MoneroRpcError) Error() string {
	res, _ := json.Marshal(e)
	return string(res)
}

/**
	JSON RPC METHODS
**/

// get_block_count
type GetBlockCountParams EmptyMoneroRpcParams
type GetBlockCountResult struct {
	Count uint64 `json:"count"`
	JsonRpcFooter
}

// on_get_block_hash
type OnGetBlockHashParams [1]uint64
type OnGetBlockHashResult string

// get_block_template
type GetBlockTemplateParams struct {
	WalletAddress string `json:"wallet_address"`
	ReserveSize   uint64 `json:"reserve_size"`
}
type GetBlockTemplateResult struct {
	BlockhashingBlob  string `json:"blockhashing_blob"`
	BlocktemplateBlob string `json:"blocktemplate_blob"`
	Difficulty        uint64 `json:"difficulty"`
	DifficultyTop64   uint64 `json:"difficulty_top64"`
	ExpectedReward    uint64 `json:"expected_reward"`
	Height            uint64 `json:"height"`
	NextSeedHash      string `json:"next_seed_hash"`
	PrevHash          string `json:"prev_hash"`
	ReservedOffset    uint64 `json:"reserved_offset"`
	SeedHash          string `json:"seed_hash"`
	SeedHeight        uint64 `json:"seed_height"`
	WideDifficulty    string `json:"wide_difficulty"`
	JsonRpcFooter
}

// get_last_block_header
type GetLastBlockHeaderParams struct {
	FillPowHash bool `json:"fill_pow_hash"`
}
type BlockHeader struct {
	BlockSize                 uint64 `json:"block_size"`
	BlockWeight               uint64 `json:"block_weight"`
	CumulativeDifficulty      uint64 `json:"cumulative_difficulty"`
	CumulativeDifficultyTop64 uint64 `json:"cumulative_difficulty_top64"`
	Depth                     uint64 `json:"depth"`
	Difficulty                uint64 `json:"difficulty"`
	DifficultyTop64           uint64 `json:"difficulty_top64"`
	Hash                      string `json:"hash"`
	Height                    uint64 `json:"height"`
	LongTermWeight            uint64 `json:"long_term_weight"`
	MajorVersion              uint   `json:"major_version"`
	MinerTxHash               string `json:"miner_tx_hash"`
	MinorVersion              uint   `json:"minor_version"`
	Nonce                     uint64 `json:"nonce"`
	NumTxes                   uint   `json:"num_txes"`
	OrphanStatus              bool   `json:"orphan_status"`
	PowHash                   string `json:"pow_hash"`
	PrevHash                  string `json:"prev_hash"`
	Reward                    uint64 `json:"reward"`
	Timestamp                 uint64 `json:"timestamp"`
	WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
	WideDifficulty            string `json:"wide_difficulty"`
}
type GetLastBlockHeaderResult struct {
	BlockHeader BlockHeader `json:"block_header"`
	Credits     uint64      `json:"credits"`
	TopHash     string      `json:"top_hash"`
	JsonRpcFooter
}

// get_block_header_by_hash
type GetBlockHeaderByHashParams struct {
	GetLastBlockHeaderParams
	Hash string `json:"hash"`
}
type GetBlockHeaderByHashResult GetLastBlockHeaderResult

// get_block_header_by_height
type GetBlockHeaderByHeightParams struct {
	GetLastBlockHeaderParams
	Height uint64 `json:"height"`
}
type GetBlockHeaderByHeightResult GetLastBlockHeaderResult

// get_block_headers_range
type GetBlockHeadersRangeParams struct {
	GetLastBlockHeaderParams
	StartHeight uint64 `json:"start_height"`
	EndHeight   uint64 `json:"end_height"`
}
type GetBlockHeadersRangeResult struct {
	Headers []BlockHeader `json:"headers"`
	Credits uint64        `json:"credits"`
	TopHash string        `json:"top_hash"`
	JsonRpcFooter
}

/**
	OTHER RPC METHODS
**/

// get_height
type GetHeightResponse struct {
	Hash   string         `json:"hash"`
	Height uint64         `json:"height"`
	Error  MoneroRpcError `json:"error"`
	JsonRpcFooter
}
