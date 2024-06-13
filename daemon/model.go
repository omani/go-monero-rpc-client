package daemon

import "encoding/json"

const (
	DEFAULT_MONERO_RPC_ENDPOINT = "/json_rpc"
)

type MoneroRpcResponse interface {
	JsonRpcResponse | OtherRpcResponse
}

type MoneroRpcRequestBody interface {
	JsonRpcRequestBody | OtherRpcRequestBody
}

type MoneroRpcRequest[T MoneroRpcRequestBody] struct {
	Endpoint string
	Body     *T
}

type MoneroRpcError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type EmptyMoneroRpcParams struct{}

func (e *MoneroRpcError) Error() string {
	res, _ := json.Marshal(e)
	return string(res)
}

/**
	JSON RPC METHODS
**/

type JsonRpcRequestParams interface {
	OnGetBlockHashParams |
		GetBlockCountParams |
		GetBlockTemplateParams |
		EmptyMoneroRpcParams |
		GetBlockHeaderDefaultParams |
		GetBlockHeaderByHashParams |
		GetBlockHeaderByHeightParams |
		GetBlockHeadersRangeParams |
		GetBlockByHashParams |
		GetBlockByHeightParams
}

type JsonRpcRequestBody interface {
	JsonRpcGenericRequestBody[OnGetBlockHashParams] |
		JsonRpcGenericRequestBody[EmptyMoneroRpcParams] |
		JsonRpcGenericRequestBody[GetBlockHeaderDefaultParams] |
		JsonRpcGenericRequestBody[GetBlockHeaderByHashParams] |
		JsonRpcGenericRequestBody[GetBlockHeaderByHeightParams] |
		JsonRpcGenericRequestBody[GetBlockHeadersRangeParams] |
		JsonRpcGenericRequestBody[GetBlockByHashParams] |
		JsonRpcGenericRequestBody[GetBlockByHeightParams] |
		JsonRpcGenericRequestBody[GetBlockTemplateParams] |
		JsonRpcGenericRequestBody[GetBlockCountParams]
}

type JsonRpcResponseResult interface {
	GetBlockCountResult |
		OnGetBlockHashResult |
		GetBlockTemplateResult |
		GetBlockHeaderResult |
		GetBlockHeadersRangeResult |
		GetBlockResult |
		GetFeeEstimateResult |
		GetVersionResult |
		GetInfoResult
}

type JsonRpcResponse interface {
	JsonRpcGenericResponse[GetBlockCountResult] |
		JsonRpcGenericResponse[OnGetBlockHashResult] |
		JsonRpcGenericResponse[GetBlockTemplateResult] |
		JsonRpcGenericResponse[GetBlockHeaderResult] |
		JsonRpcGenericResponse[GetBlockHeadersRangeResult] |
		JsonRpcGenericResponse[GetBlockResult] |
		JsonRpcGenericResponse[GetFeeEstimateResult] |
		JsonRpcGenericResponse[GetVersionResult] |
		JsonRpcGenericResponse[GetInfoResult]
}

type JsonRpcHeader struct {
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}
type JsonRpcFooter struct {
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

type JsonRpcGenericRequestBody[T JsonRpcRequestParams] struct {
	JsonRpcHeader
	Method string `json:"method"`
	Params T      `json:"params"`
}

type JsonRpcGenericResponse[T JsonRpcResponseResult] struct {
	JsonRpcHeader
	Result T              `json:"result"`
	Error  MoneroRpcError `json:"error"`
}

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
type GetBlockHeaderDefaultParams struct {
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
	Timestamp                 uint32 `json:"timestamp"`
	WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
	WideDifficulty            string `json:"wide_difficulty"`
}
type GetBlockHeaderResult struct {
	BlockHeader BlockHeader `json:"block_header"`
	Credits     uint64      `json:"credits"`
	TopHash     string      `json:"top_hash"`
	JsonRpcFooter
}

// get_block_header_by_hash
type GetBlockHeaderByHashParams struct {
	GetBlockHeaderDefaultParams
	Hash string `json:"hash"`
}

// get_block_header_by_height
type GetBlockHeaderByHeightParams struct {
	GetBlockHeaderDefaultParams
	Height uint64 `json:"height"`
}

// get_block_headers_range
type GetBlockHeadersRangeParams struct {
	GetBlockHeaderDefaultParams
	StartHeight uint64 `json:"start_height"`
	EndHeight   uint64 `json:"end_height"`
}
type GetBlockHeadersRangeResult struct {
	Headers []BlockHeader `json:"headers"`
	Credits uint64        `json:"credits"`
	TopHash string        `json:"top_hash"`
	JsonRpcFooter
}

// get_block
type GetBlockByHashParams GetBlockHeaderByHashParams
type GetBlockByHeightParams GetBlockHeaderByHeightParams

type Gen struct {
	Height uint64 `json:"height"`
}
type Vin1 struct {
	Gen Gen `json:"gen"`
}
type TaggedKey struct {
	Key     string `json:"key"`
	ViewTag string `json:"view_tag"`
}
type Target struct {
	TaggedKey TaggedKey `json:"tagged_key"`
}
type Vout1 struct {
	Amount uint64 `json:"amount"`
	Target Target `json:"target"`
}
type RctSignatures struct {
	Type uint32 `json:"type"`
}
type MinerTx struct {
	Version       uint32        `json:"version"`
	UnlockTime    uint64        `json:"unlock_time"`
	Vin           []Vin1        `json:"vin"`
	Vout          []Vout1       `json:"vout"`
	Extra         []int32       `json:"extra"`
	RctSignatures RctSignatures `json:"rct_signatures"`
}
type BlockDetails struct {
	MajorVersion uint     `json:"major_version"`
	MinorVersion uint     `json:"minor_version"`
	Timestamp    uint32   `json:"timestamp"`
	PrevId       string   `json:"prev_id"`
	Nonce        uint64   `json:"nonce"`
	MinerTx      MinerTx  `json:"miner_tx"`
	TxHashes     []string `json:"tx_hashes"`
}

type GetBlockResult struct {
	Blob         string `json:"blob"`
	MinerTxHash  string `json:"miner_tx_hash"`
	Json         string `json:"json"`
	BlockDetails BlockDetails
	GetBlockHeaderResult
}

// get_fee_estimate
type GetFeeEstimateResult struct {
	Credits          uint64   `json:"credits"`
	Fee              uint64   `json:"fee"`
	Fees             []uint64 `json:"fees"`
	QuantizationMask uint64   `json:"quantization_mask"`
	TopHash          string   `json:"top_hash"`
	JsonRpcFooter
}

// get_version
type GetVersionResult struct {
	Release bool   `json:"release"`
	Version uint32 `json:"version"`
	JsonRpcFooter
}

// get_info
type GetInfoResult struct {
	AdjustedTime              uint64 `json:"adjusted_time"`
	AltBlocksCount            uint32 `json:"alt_blocks_count"`
	BlockSizeLimit            uint64 `json:"block_size_limit"`
	BlockSizeMedian           uint64 `json:"block_size_median"`
	BlockWeightLimit          uint64 `json:"block_weight_limit"`
	BlockWeightMedian         uint64 `json:"block_weight_median"`
	BootstrapDaemonAddress    string `json:"bootstrap_daemon_address"`
	BusySyncing               bool   `json:"busy_syncing"`
	Credits                   uint64 `json:"credits"`
	CumulativeDifficulty      uint64 `json:"cumulative_difficulty"`
	CumulativeDifficultyTop64 uint64 `json:"cumulative_difficulty_top64"`
	DatabaseSize              uint64 `json:"database_size"`
	Difficulty                uint64 `json:"difficulty"`
	DifficultyTop64           uint64 `json:"difficulty_top64"`
	FreeSpace                 uint64 `json:"free_space"`
	GreyPeerlistSize          uint32 `json:"grey_peerlist_size"`
	Height                    uint64 `json:"height"`
	HeightWithoutBootstrap    uint64 `json:"height_without_bootstrap"`
	IncomingConnectionsCount  uint32 `json:"incoming_connections_count"`
	Mainnet                   bool   `json:"mainnet"`
	Nettype                   string `json:"nettype"`
	Offline                   bool   `json:"offline"`
	OutgoingConnectionsCount  uint32 `json:"outgoing_connections_count"`
	RpcConnectionsCount       uint32 `json:"rpc_connections_count"`
	Stagenet                  bool   `json:"stagenet"`
	StartTime                 uint64 `json:"start_time"`
	Synchronized              bool   `json:"synchronized"`
	Target                    uint32 `json:"target"`
	TargetHeight              uint64 `json:"target_height"`
	Testnet                   bool   `json:"testnet"`
	TopBlockHash              string `json:"top_block_hash"`
	TopHash                   string `json:"top_hash"`
	TxCount                   uint64 `json:"tx_count"`
	TxPoolSize                uint32 `json:"tx_pool_size"`
	UpdateAvailable           bool   `json:"update_available"`
	Version                   string `json:"version"`
	WasBootstrapEverUsed      bool   `json:"was_bootstrap_ever_used"`
	WhitePeerlistSize         uint32 `json:"white_peerlist_size"`
	WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
	WideDifficulty            string `json:"wide_difficulty"`
	JsonRpcFooter
}

/**
	OTHER RPC METHODS
**/

type OtherRpcRequestBody interface {
	EmptyMoneroRpcParams |
		GetTransactionsParams
}

type OtherRpcResponse interface {
	GetHeightResponse |
		GetTransactionPoolResponse |
		GetTransactionsResponse
}

// get_height
type GetHeightResponse struct {
	Hash   string         `json:"hash"`
	Height uint64         `json:"height"`
	Error  MoneroRpcError `json:"error"`
	JsonRpcFooter
}

// get_transaction_pool
type SpentKeyImage struct {
	IdHash    string   `json:"id_hash"`
	TxsHashes []string `json:"txs_hashes"`
}
type Key struct {
	Amount     uint64  `json:"amount"`
	KeyOffsets []int64 `json:"key_offsets"`
	KeyImage   string  `json:"k_image"`
}
type Vin2 struct {
	Key Key `json:"key"`
}
type EcdhInfo struct {
	Amount      string `json:"amount"`
	TruncAmount string `json:"trunc_amount"`
}
type RctSignature struct {
	Type     int32      `type:"txnFee"`
	TxnFee   uint64     `json:"txnFee"`
	EcdhInfo []EcdhInfo `json:"ecdhInfo"`
	OutPk    []string   `json:"outPk"`
}
type CLSAG struct {
	D  string   `json:"D"`
	C1 string   `json:"c1"`
	S  []string `json:"s"`
}
type Bpp struct {
	A  string   `json:"A"`
	A1 string   `json:"A1"`
	B  string   `json:"B"`
	L  []string `json:"L"`
	R  []string `json:"R"`
	R1 string   `json:"r1"`
	D1 string   `json:"d1"`
	S1 string   `json:"s1"`
}
type RctsigPrunable struct {
	CLSAGs     []CLSAG  `json:"CLSAGs"`
	Bpp        []Bpp    `json:"bpp"`
	Nbp        int32    `json:"nbp"`
	PseudoOuts []string `json:"pseudoOuts"`
}
type MoneroTxInfo struct {
	Version        uint32         `json:"version"`
	UnlockTime     uint64         `json:"unlock_time"`
	Vin            []Vin2         `json:"vin"`
	Vout           []Vout1        `json:"vout"`
	Extra          []int32        `json:"extra"`
	RctSignatures  RctSignature   `json:"rct_signatures"`
	RctsigPrunable RctsigPrunable `json:"rctsig_prunable"`
}
type MoneroTx struct {
	BlobSize           uint64 `json:"blob_size"`
	DoNotRelay         bool   `json:"do_not_relay"`
	DoubleSpendSeen    bool   `json:"double_spend_seen"`
	Fee                uint64 `json:"fee"`
	IdHash             string `json:"id_hash"`
	KeptByBlock        bool   `json:"kept_by_block"`
	LastFailedHeight   uint64 `json:"last_failed_height"`
	LastFailedIdHash   string `json:"last_failed_id_hash"`
	LastRelayedTime    uint64 `json:"last_relayed_time"`
	MaxUsedBlockHeight uint64 `json:"max_used_block_height"`
	MaxUsedBlockIdHash string `json:"max_used_block_id_hash"`
	ReceiveTime        uint64 `json:"receive_time"`
	Relayed            bool   `json:"relayed"`
	TxBlob             string `json:"tx_blob"`
	TxJson             string `json:"tx_json"`
	Weight             uint64 `json:"weight"`
	TxInfo             MoneroTxInfo
}
type GetTransactionPoolResponse struct {
	Credits        uint64          `json:"credits"`
	SpentKeyImages []SpentKeyImage `json:"spent_key_images"`
	TopHash        string          `json:"top_hash"`
	Transactions   []MoneroTx      `json:"transactions"`
	Error          MoneroRpcError  `json:"error"`
	JsonRpcFooter
}

// get_transactions
type GetTransactionsParams struct {
	TxHashes     []string `json:"txs_hashes"`
	DecodeAsJson bool     `json:"decode_as_json"`
	Prune        bool     `json:"prune"`
	Split        bool     `json:"split"`
}

type MoneroTx1 struct {
	AsHex           string   `json:"as_hex"`
	AsJson          string   `json:"as_json"`
	BlockHeight     uint64   `json:"block_height"`
	BlockTimestamp  uint64   `json:"block_timestamp"`
	Confirmations   uint64   `json:"confirmations"`
	DoubleSpendSeen bool     `json:"double_spend_seen"`
	InPool          bool     `json:"in_pool"`
	OutputIndices   []uint64 `json:"output_indices"`
	PrunableAsHex   string   `json:"prunable_as_hex"`
	PrunableHash    string   `json:"prunable_hash"`
	PrunedAsHex     string   `json:"pruned_as_hex"`
	TxHash          string   `json:"tx_hash"`
	TxInfo          MoneroTxInfo
}
type GetTransactionsResponse struct {
	Credits  uint64         `json:"credits"`
	MissedTx []string       `json:"missed_tx"`
	TopHash  string         `json:"top_hash"`
	Txs      []MoneroTx1    `json:"txs"`
	Error    MoneroRpcError `json:"error"`
	JsonRpcFooter
}
