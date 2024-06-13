package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/monero-ecosystem/go-monero-rpc-client/daemon"
	"github.com/stretchr/testify/assert"
)

var defaultMoneroRpcHeader = daemon.JsonRpcHeader{Id: "0", Jsonrpc: "2.0"}
var defaultMoneroRpcFooter = daemon.JsonRpcFooter{Status: "OK", Untrusted: false}

func createTestDaemonRpcClient(u string) (daemon.IDaemonRpcClient, error) {
	u1, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return daemon.NewDaemonRpcClient(daemon.NewRpcConnection(u1, "", "")), nil
}

func compareRequestBody[T daemon.MoneroRpcRequestBody](r *http.Request, expected *T) bool {
	if expected == nil {
		return true
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	var actual T
	if err := json.Unmarshal(data, &actual); err != nil {
		return false
	}

	return reflect.DeepEqual(&actual, expected)
}
func daemonRpcTestServerCheck[B daemon.MoneroRpcRequestBody](r *http.Request, expected *daemon.MoneroRpcRequest[B]) bool {
	if r == nil {
		return false
	}

	if r.URL.Path != expected.Endpoint ||
		r.Method != http.MethodPost ||
		r.Header.Get("Accept") != "application/json" {
		return false
	}

	return compareRequestBody(r, expected.Body)
}
func getDaemonRpcTestServer[B daemon.MoneroRpcRequestBody](req *daemon.MoneroRpcRequest[B], res *string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if daemonRpcTestServerCheck(r, req) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(*res))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
}

func TestGetCurrentHeight(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.EmptyMoneroRpcParams]{Endpoint: "/get_height", Body: nil}
	exres := `{
				"hash": "7e23a28cfa6df925d5b63940baf60b83c0cbb65da95f49b19e7cf0ce7dd709ce",
				"height": 2287217,
				"status": "OK",
				"untrusted": false
			}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.GetHeightResponse{
		Hash:          "7e23a28cfa6df925d5b63940baf60b83c0cbb65da95f49b19e7cf0ce7dd709ce",
		Height:        2287217,
		Error:         daemon.MoneroRpcError{},
		JsonRpcFooter: defaultMoneroRpcFooter}
	actual, err := test_daemon.GetCurrentHeight()
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockCount(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockCountParams]{JsonRpcHeader: defaultMoneroRpcHeader, Method: "get_block_count", Params: daemon.GetBlockCountParams{}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockCountParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
		  "count": 993163,
		  "status": "OK",
		  "untrusted": false
		}
	  }`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockCountResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockCountResult{
			Count:         993163,
			JsonRpcFooter: defaultMoneroRpcFooter,
		},
		Error: daemon.MoneroRpcError{},
	}
	actual, err := test_daemon.GetBlockCount()
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestOnGetBlockHash(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.OnGetBlockHashParams]{JsonRpcHeader: defaultMoneroRpcHeader, Method: "on_get_block_hash", Params: [1]uint64{912345}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.OnGetBlockHashParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6"
		}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.OnGetBlockHashResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result:        "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
		Error:         daemon.MoneroRpcError{}}

	actual, err := test_daemon.OnGetBlockHash(exreq.Body.Params[0])
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockTemplate(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockTemplateParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_block_template",
		Params: daemon.GetBlockTemplateParams{
			WalletAddress: "44GBHzv6ZyQdJkjqZje6KLZ3xSyN1hBSFAnLP6EAqJtCRVzMzZmeXTC2AHKDS9aEDTRKmo6a6o9r9j86pYfhCWDkKjbtcns", ReserveSize: 60}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockTemplateParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
		  "blockhashing_blob": "0e0ed286da8006ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd00000000d130d22cf308b308498bbc16e2e955e7dbd691e6a8fab805f98ad82e6faa8bcc06",
		  "blocktemplate_blob": "0e0ed286da8006ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd0000000002abc78b0101ffefc68b0101fcfcf0d4b422025014bb4a1eade6622fd781cb1063381cad396efa69719b41aa28b4fce8c7ad4b5f019ce1dc670456b24a5e03c2d9058a2df10fec779e2579753b1847b74ee644f16b023c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000051399a1bc46a846474f5b33db24eae173a26393b976054ee14f9feefe99925233802867097564c9db7a36af5bb5ed33ab46e63092bd8d32cef121608c3258edd55562812e21cc7e3ac73045745a72f7d74581d9a0849d6f30e8b2923171253e864f4e9ddea3acb5bc755f1c4a878130a70c26297540bc0b7a57affb6b35c1f03d8dbd54ece8457531f8cba15bb74516779c01193e212050423020e45aa2c15dcb",
		  "difficulty": 226807339040,
		  "difficulty_top64": 0,
		  "expected_reward": 1182367759996,
		  "height": 2286447,
		  "next_seed_hash": "",
		  "prev_hash": "ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd",
		  "reserved_offset": 130,
		  "seed_hash": "d432f499205150873b2572b5f033c9c6e4b7c6f3394bd2dd93822cd7085e7307",
		  "seed_height": 2285568,
		  "status": "OK",
		  "untrusted": false,
		  "wide_difficulty": "0x34cec55820"
		}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockTemplateResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockTemplateResult{
			BlockhashingBlob:  "0e0ed286da8006ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd00000000d130d22cf308b308498bbc16e2e955e7dbd691e6a8fab805f98ad82e6faa8bcc06",
			BlocktemplateBlob: "0e0ed286da8006ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd0000000002abc78b0101ffefc68b0101fcfcf0d4b422025014bb4a1eade6622fd781cb1063381cad396efa69719b41aa28b4fce8c7ad4b5f019ce1dc670456b24a5e03c2d9058a2df10fec779e2579753b1847b74ee644f16b023c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000051399a1bc46a846474f5b33db24eae173a26393b976054ee14f9feefe99925233802867097564c9db7a36af5bb5ed33ab46e63092bd8d32cef121608c3258edd55562812e21cc7e3ac73045745a72f7d74581d9a0849d6f30e8b2923171253e864f4e9ddea3acb5bc755f1c4a878130a70c26297540bc0b7a57affb6b35c1f03d8dbd54ece8457531f8cba15bb74516779c01193e212050423020e45aa2c15dcb",
			Difficulty:        226807339040,
			DifficultyTop64:   0,
			ExpectedReward:    1182367759996,
			Height:            2286447,
			NextSeedHash:      "",
			PrevHash:          "ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd",
			ReservedOffset:    130,
			SeedHash:          "d432f499205150873b2572b5f033c9c6e4b7c6f3394bd2dd93822cd7085e7307",
			SeedHeight:        2285568,
			WideDifficulty:    "0x34cec55820",
			JsonRpcFooter:     defaultMoneroRpcFooter},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockTemplate(exreq.Body.Params.WalletAddress, exreq.Body.Params.ReserveSize)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetLastBlockHeader(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeaderDefaultParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_last_block_header",
		Params:        daemon.GetBlockHeaderDefaultParams{}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeaderDefaultParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
		  "block_header": {
			"block_size": 5500,
			"block_weight": 5500,
			"cumulative_difficulty": 86164894009456483,
			"cumulative_difficulty_top64": 0,
			"depth": 0,
			"difficulty": 227026389695,
			"difficulty_top64": 0,
			"hash": "a6ad87cf357a1aac1ee1d7cb0afa4c2e653b0b1ab7d5bf6af310333e43c59dd0",
			"height": 2286454,
			"long_term_weight": 5500,
			"major_version": 14,
			"miner_tx_hash": "a474f87de1645ff14c5e90c477b07f9bc86a22fb42909caa0705239298da96d0",
			"minor_version": 14,
			"nonce": 249602367,
			"num_txes": 3,
			"orphan_status": false,
			"pow_hash": "",
			"prev_hash": "fa17fefe1d05da775a61a3dc33d9e199d12af167ef0ab37e52b51e8487b50f25",
			"reward": 1181337498013,
			"timestamp": 1612088597,
			"wide_cumulative_difficulty": "0x1321e83bb8af763",
			"wide_difficulty": "0x34dbd3cabf"
		  },
		  "credits": 0,
		  "status": "OK",
		  "top_hash": "",
		  "untrusted": false
		}
	  }`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockHeaderResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockHeaderResult{
			BlockHeader: daemon.BlockHeader{
				BlockSize:                 5500,
				BlockWeight:               5500,
				CumulativeDifficulty:      86164894009456483,
				CumulativeDifficultyTop64: 0,
				Depth:                     0,
				Difficulty:                227026389695,
				DifficultyTop64:           0,
				Hash:                      "a6ad87cf357a1aac1ee1d7cb0afa4c2e653b0b1ab7d5bf6af310333e43c59dd0",
				Height:                    2286454,
				LongTermWeight:            5500,
				MajorVersion:              14,
				MinerTxHash:               "a474f87de1645ff14c5e90c477b07f9bc86a22fb42909caa0705239298da96d0",
				MinorVersion:              14,
				Nonce:                     249602367,
				NumTxes:                   3,
				OrphanStatus:              false,
				PowHash:                   "",
				PrevHash:                  "fa17fefe1d05da775a61a3dc33d9e199d12af167ef0ab37e52b51e8487b50f25",
				Reward:                    1181337498013,
				Timestamp:                 1612088597,
				WideCumulativeDifficulty:  "0x1321e83bb8af763",
				WideDifficulty:            "0x34dbd3cabf"},
			Credits:       0,
			TopHash:       "",
			JsonRpcFooter: defaultMoneroRpcFooter},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetLastBlockHeader(false)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetLastBlockHeaderByHash(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeaderByHashParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_block_header_by_hash",
		Params: daemon.GetBlockHeaderByHashParams{
			GetBlockHeaderDefaultParams: daemon.GetBlockHeaderDefaultParams{FillPowHash: false},
			Hash:                        "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6"}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeaderByHashParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
		  "block_header": {
			"block_size": 210,
			"block_weight": 210,
			"cumulative_difficulty": 754734824984346,
			"cumulative_difficulty_top64": 0,
			"depth": 1374113,
			"difficulty": 815625611,
			"difficulty_top64": 0,
			"hash": "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
			"height": 912345,
			"long_term_weight": 210,
			"major_version": 1,
			"miner_tx_hash": "c7da3965f25c19b8eb7dd8db48dcd4e7c885e2491db77e289f0609bf8e08ec30",
			"minor_version": 2,
			"nonce": 1646,
			"num_txes": 0,
			"orphan_status": false,
			"pow_hash": "",
			"prev_hash": "b61c58b2e0be53fad5ef9d9731a55e8a81d972b8d90ed07c04fd37ca6403ff78",
			"reward": 7388968946286,
			"timestamp": 1452793716,
			"wide_cumulative_difficulty": "0x2ae6d65248f1a",
			"wide_difficulty": "0x309d758b"
		  },
		  "credits": 0,
		  "status": "OK",
		  "top_hash": "",
		  "untrusted": false
		}
	  }`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockHeaderResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockHeaderResult{
			BlockHeader: daemon.BlockHeader{
				BlockSize:                 210,
				BlockWeight:               210,
				CumulativeDifficulty:      754734824984346,
				CumulativeDifficultyTop64: 0,
				Depth:                     1374113,
				Difficulty:                815625611,
				DifficultyTop64:           0,
				Hash:                      "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
				Height:                    912345,
				LongTermWeight:            210,
				MajorVersion:              1,
				MinerTxHash:               "c7da3965f25c19b8eb7dd8db48dcd4e7c885e2491db77e289f0609bf8e08ec30",
				MinorVersion:              2,
				Nonce:                     1646,
				NumTxes:                   0,
				OrphanStatus:              false,
				PowHash:                   "",
				PrevHash:                  "b61c58b2e0be53fad5ef9d9731a55e8a81d972b8d90ed07c04fd37ca6403ff78",
				Reward:                    7388968946286,
				Timestamp:                 1452793716,
				WideCumulativeDifficulty:  "0x2ae6d65248f1a",
				WideDifficulty:            "0x309d758b"},
			Credits:       0,
			TopHash:       "",
			JsonRpcFooter: defaultMoneroRpcFooter},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockHeaderByHash(exreq.Body.Params.FillPowHash, exreq.Body.Params.Hash)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetLastBlockHeaderByHeight(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeaderByHeightParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_block_header_by_height",
		Params: daemon.GetBlockHeaderByHeightParams{
			GetBlockHeaderDefaultParams: daemon.GetBlockHeaderDefaultParams{FillPowHash: false},
			Height:                      912345}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeaderByHeightParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
				"id": "0",
				"jsonrpc": "2.0",
				"result": {
					"block_header": {
					"block_size": 210,
					"block_weight": 210,
					"cumulative_difficulty": 754734824984346,
					"cumulative_difficulty_top64": 0,
					"depth": 1374118,
					"difficulty": 815625611,
					"difficulty_top64": 0,
					"hash": "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
					"height": 912345,
					"long_term_weight": 210,
					"major_version": 1,
					"miner_tx_hash": "c7da3965f25c19b8eb7dd8db48dcd4e7c885e2491db77e289f0609bf8e08ec30",
					"minor_version": 2,
					"nonce": 1646,
					"num_txes": 0,
					"orphan_status": false,
					"pow_hash": "",
					"prev_hash": "b61c58b2e0be53fad5ef9d9731a55e8a81d972b8d90ed07c04fd37ca6403ff78",
					"reward": 7388968946286,
					"timestamp": 1452793716,
					"wide_cumulative_difficulty": "0x2ae6d65248f1a",
					"wide_difficulty": "0x309d758b"
					},
					"credits": 0,
					"status": "OK",
					"top_hash": "",
					"untrusted": false
				}
			}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockHeaderResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockHeaderResult{
			BlockHeader: daemon.BlockHeader{
				BlockSize:                 210,
				BlockWeight:               210,
				CumulativeDifficulty:      754734824984346,
				CumulativeDifficultyTop64: 0,
				Depth:                     1374118,
				Difficulty:                815625611,
				DifficultyTop64:           0,
				Hash:                      "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
				Height:                    912345,
				LongTermWeight:            210,
				MajorVersion:              1,
				MinerTxHash:               "c7da3965f25c19b8eb7dd8db48dcd4e7c885e2491db77e289f0609bf8e08ec30",
				MinorVersion:              2,
				Nonce:                     1646,
				NumTxes:                   0,
				OrphanStatus:              false,
				PowHash:                   "",
				PrevHash:                  "b61c58b2e0be53fad5ef9d9731a55e8a81d972b8d90ed07c04fd37ca6403ff78",
				Reward:                    7388968946286,
				Timestamp:                 1452793716,
				WideCumulativeDifficulty:  "0x2ae6d65248f1a",
				WideDifficulty:            "0x309d758b"},
			Credits:       0,
			TopHash:       "",
			JsonRpcFooter: defaultMoneroRpcFooter},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockHeaderByHeight(exreq.Body.Params.FillPowHash, exreq.Body.Params.Height)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockHeadersRange(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeadersRangeParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_block_headers_range",
		Params: daemon.GetBlockHeadersRangeParams{
			GetBlockHeaderDefaultParams: daemon.GetBlockHeaderDefaultParams{FillPowHash: false},
			StartHeight:                 1545999,
			EndHeight:                   1546000}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockHeadersRangeParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
			"credits": 0,
			"headers": [{
			"block_size": 301413,
			"block_weight": 301413,
			"cumulative_difficulty": 13185267971483472,
			"cumulative_difficulty_top64": 0,
			"depth": 740464,
			"difficulty": 134636057921,
			"difficulty_top64": 0,
			"hash": "86d1d20a40cefcf3dd410ff6967e0491613b77bf73ea8f1bf2e335cf9cf7d57a",
			"height": 1545999,
			"long_term_weight": 301413,
			"major_version": 6,
			"miner_tx_hash": "9909c6f8a5267f043c3b2b079fb4eacc49ef9c1dee1c028eeb1a259b95e6e1d9",
			"minor_version": 6,
			"nonce": 3246403956,
			"num_txes": 20,
			"orphan_status": false,
			"pow_hash": "",
			"prev_hash": "0ef6e948f77b8f8806621003f5de24b1bcbea150bc0e376835aea099674a5db5",
			"reward": 5025593029981,
			"timestamp": 1523002893,
			"wide_cumulative_difficulty": "0x2ed7ee6db56750",
			"wide_difficulty": "0x1f58ef3541"
			},{
			"block_size": 13322,
			"block_weight": 13322,
			"cumulative_difficulty": 13185402687569710,
			"cumulative_difficulty_top64": 0,
			"depth": 740463,
			"difficulty": 134716086238,
			"difficulty_top64": 0,
			"hash": "b408bf4cfcd7de13e7e370c84b8314c85b24f0ba4093ca1d6eeb30b35e34e91a",
			"height": 1546000,
			"long_term_weight": 13322,
			"major_version": 7,
			"miner_tx_hash": "7f749c7c64acb35ef427c7454c45e6688781fbead9bbf222cb12ad1a96a4e8f6",
			"minor_version": 7,
			"nonce": 3737164176,
			"num_txes": 1,
			"orphan_status": false,
			"pow_hash": "",
			"prev_hash": "86d1d20a40cefcf3dd410ff6967e0491613b77bf73ea8f1bf2e335cf9cf7d57a",
			"reward": 4851952181070,
			"timestamp": 1523002931,
			"wide_cumulative_difficulty": "0x2ed80dcb69bf2e",
			"wide_difficulty": "0x1f5db457de"
			}],
			"status": "OK",
			"top_hash": "",
			"untrusted": false
		}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	headers := []daemon.BlockHeader{
		{
			BlockSize:                 301413,
			BlockWeight:               301413,
			CumulativeDifficulty:      13185267971483472,
			CumulativeDifficultyTop64: 0,
			Depth:                     740464,
			Difficulty:                134636057921,
			DifficultyTop64:           0,
			Hash:                      "86d1d20a40cefcf3dd410ff6967e0491613b77bf73ea8f1bf2e335cf9cf7d57a",
			Height:                    1545999,
			LongTermWeight:            301413,
			MajorVersion:              6,
			MinerTxHash:               "9909c6f8a5267f043c3b2b079fb4eacc49ef9c1dee1c028eeb1a259b95e6e1d9",
			MinorVersion:              6,
			Nonce:                     3246403956,
			NumTxes:                   20,
			OrphanStatus:              false,
			PowHash:                   "",
			PrevHash:                  "0ef6e948f77b8f8806621003f5de24b1bcbea150bc0e376835aea099674a5db5",
			Reward:                    5025593029981,
			Timestamp:                 1523002893,
			WideCumulativeDifficulty:  "0x2ed7ee6db56750",
			WideDifficulty:            "0x1f58ef3541"},
		{
			BlockSize:                 13322,
			BlockWeight:               13322,
			CumulativeDifficulty:      13185402687569710,
			CumulativeDifficultyTop64: 0,
			Depth:                     740463,
			Difficulty:                134716086238,
			DifficultyTop64:           0,
			Hash:                      "b408bf4cfcd7de13e7e370c84b8314c85b24f0ba4093ca1d6eeb30b35e34e91a",
			Height:                    1546000,
			LongTermWeight:            13322,
			MajorVersion:              7,
			MinerTxHash:               "7f749c7c64acb35ef427c7454c45e6688781fbead9bbf222cb12ad1a96a4e8f6",
			MinorVersion:              7,
			Nonce:                     3737164176,
			NumTxes:                   1,
			OrphanStatus:              false,
			PowHash:                   "",
			PrevHash:                  "86d1d20a40cefcf3dd410ff6967e0491613b77bf73ea8f1bf2e335cf9cf7d57a",
			Reward:                    4851952181070,
			Timestamp:                 1523002931,
			WideCumulativeDifficulty:  "0x2ed80dcb69bf2e",
			WideDifficulty:            "0x1f5db457de"}}
	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockHeadersRangeResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockHeadersRangeResult{
			Headers:       headers,
			Credits:       0,
			TopHash:       "",
			JsonRpcFooter: defaultMoneroRpcFooter},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockHeadersRange(exreq.Body.Params.FillPowHash, exreq.Body.Params.StartHeight, exreq.Body.Params.EndHeight)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockByHeight(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockByHeightParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_block",
		Params: daemon.GetBlockByHeightParams{
			GetBlockHeaderDefaultParams: daemon.GetBlockHeaderDefaultParams{FillPowHash: false},
			Height:                      2751506}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockByHeightParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
			"blob": "1010c58bab9b06b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7807e07f502cef8a70101ff92f8a7010180e0a596bb1103d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85d034019f629d8b36bd16a2bfce3ea80c31dc4d8762c67165aec21845494e32b7582fe00211000000297a787a000000000000000000000000",
			"block_header": {
			"block_size": 106,
			"block_weight": 106,
			"cumulative_difficulty": 236046001376524168,
			"cumulative_difficulty_top64": 0,
			"depth": 40,
			"difficulty": 313732272488,
			"difficulty_top64": 0,
			"hash": "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428",
			"height": 2751506,
			"long_term_weight": 176470,
			"major_version": 16,
			"miner_tx_hash": "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
			"minor_version": 16,
			"nonce": 4110909056,
			"num_txes": 0,
			"orphan_status": false,
			"pow_hash": "",
			"prev_hash": "b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7",
			"reward": 600000000000,
			"timestamp": 1667941829,
			"wide_cumulative_difficulty": "0x3469a966eb2f788",
			"wide_difficulty": "0x490be69168"
			},
			"credits": 0,
			"json": "{\n  \"major_version\": 16, \n  \"minor_version\": 16, \n  \"timestamp\": 1667941829, \n  \"prev_id\": \"b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7\", \n  \"nonce\": 4110909056, \n  \"miner_tx\": {\n    \"version\": 2, \n    \"unlock_time\": 2751566, \n    \"vin\": [ {\n        \"gen\": {\n          \"height\": 2751506\n        }\n      }\n    ], \n    \"vout\": [ {\n        \"amount\": 600000000000, \n        \"target\": {\n          \"tagged_key\": {\n            \"key\": \"d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85\", \n            \"view_tag\": \"d0\"\n          }\n        }\n      }\n    ], \n    \"extra\": [ 1, 159, 98, 157, 139, 54, 189, 22, 162, 191, 206, 62, 168, 12, 49, 220, 77, 135, 98, 198, 113, 101, 174, 194, 24, 69, 73, 78, 50, 183, 88, 47, 224, 2, 17, 0, 0, 0, 41, 122, 120, 122, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0\n    ], \n    \"rct_signatures\": {\n      \"type\": 0\n    }\n  }, \n  \"tx_hashes\": [ ]\n}",
			"miner_tx_hash": "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
			"status": "OK",
			"top_hash": "",
			"untrusted": false
		}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockResult{
			Blob: "1010c58bab9b06b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7807e07f502cef8a70101ff92f8a7010180e0a596bb1103d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85d034019f629d8b36bd16a2bfce3ea80c31dc4d8762c67165aec21845494e32b7582fe00211000000297a787a000000000000000000000000",
			GetBlockHeaderResult: daemon.GetBlockHeaderResult{
				BlockHeader: daemon.BlockHeader{
					BlockSize:                 106,
					BlockWeight:               106,
					CumulativeDifficulty:      236046001376524168,
					CumulativeDifficultyTop64: 0,
					Depth:                     40,
					Difficulty:                313732272488,
					DifficultyTop64:           0,
					Hash:                      "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428",
					Height:                    2751506,
					LongTermWeight:            176470,
					MajorVersion:              16,
					MinerTxHash:               "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
					MinorVersion:              16,
					Nonce:                     4110909056,
					NumTxes:                   0,
					OrphanStatus:              false,
					PowHash:                   "",
					PrevHash:                  "b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7",
					Reward:                    600000000000,
					Timestamp:                 1667941829,
					WideCumulativeDifficulty:  "0x3469a966eb2f788",
					WideDifficulty:            "0x490be69168"},
				Credits:       0,
				TopHash:       "",
				JsonRpcFooter: defaultMoneroRpcFooter,
			},
			Json: "{\n  \"major_version\": 16, \n  \"minor_version\": 16, \n  \"timestamp\": 1667941829, \n  \"prev_id\": \"b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7\", \n  \"nonce\": 4110909056, \n  \"miner_tx\": {\n    \"version\": 2, \n    \"unlock_time\": 2751566, \n    \"vin\": [ {\n        \"gen\": {\n          \"height\": 2751506\n        }\n      }\n    ], \n    \"vout\": [ {\n        \"amount\": 600000000000, \n        \"target\": {\n          \"tagged_key\": {\n            \"key\": \"d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85\", \n            \"view_tag\": \"d0\"\n          }\n        }\n      }\n    ], \n    \"extra\": [ 1, 159, 98, 157, 139, 54, 189, 22, 162, 191, 206, 62, 168, 12, 49, 220, 77, 135, 98, 198, 113, 101, 174, 194, 24, 69, 73, 78, 50, 183, 88, 47, 224, 2, 17, 0, 0, 0, 41, 122, 120, 122, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0\n    ], \n    \"rct_signatures\": {\n      \"type\": 0\n    }\n  }, \n  \"tx_hashes\": [ ]\n}",
			BlockDetails: daemon.BlockDetails{
				MajorVersion: 16,
				MinorVersion: 16,
				Timestamp:    1667941829,
				PrevId:       "b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7",
				Nonce:        4110909056,
				MinerTx: daemon.MinerTx{
					Version:    2,
					UnlockTime: 2751566,
					Vin: []daemon.Vin1{
						{
							Gen: daemon.Gen{Height: 2751506},
						},
					},
					Vout: []daemon.Vout1{
						{
							Amount: 600000000000,
							Target: daemon.Target{
								TaggedKey: daemon.TaggedKey{
									Key:     "d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85",
									ViewTag: "d0",
								},
							},
						},
					},
					Extra:         []int32{1, 159, 98, 157, 139, 54, 189, 22, 162, 191, 206, 62, 168, 12, 49, 220, 77, 135, 98, 198, 113, 101, 174, 194, 24, 69, 73, 78, 50, 183, 88, 47, 224, 2, 17, 0, 0, 0, 41, 122, 120, 122, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					RctSignatures: daemon.RctSignatures{Type: 0},
				},
				TxHashes: []string{},
			},
			MinerTxHash: "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
		},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockByHeight(exreq.Body.Params.FillPowHash, exreq.Body.Params.Height)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockByHash(t *testing.T) {
	exreqBody := &daemon.JsonRpcGenericRequestBody[daemon.GetBlockByHashParams]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Method:        "get_block",
		Params: daemon.GetBlockByHashParams{
			GetBlockHeaderDefaultParams: daemon.GetBlockHeaderDefaultParams{FillPowHash: false},
			Hash:                        "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428"}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.GetBlockByHashParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: exreqBody}
	exres := `{
		"id": "0",
		"jsonrpc": "2.0",
		"result": {
			"blob": "1010c58bab9b06b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7807e07f502cef8a70101ff92f8a7010180e0a596bb1103d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85d034019f629d8b36bd16a2bfce3ea80c31dc4d8762c67165aec21845494e32b7582fe00211000000297a787a000000000000000000000000",
			"block_header": {
			"block_size": 106,
			"block_weight": 106,
			"cumulative_difficulty": 236046001376524168,
			"cumulative_difficulty_top64": 0,
			"depth": 40,
			"difficulty": 313732272488,
			"difficulty_top64": 0,
			"hash": "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428",
			"height": 2751506,
			"long_term_weight": 176470,
			"major_version": 16,
			"miner_tx_hash": "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
			"minor_version": 16,
			"nonce": 4110909056,
			"num_txes": 0,
			"orphan_status": false,
			"pow_hash": "",
			"prev_hash": "b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7",
			"reward": 600000000000,
			"timestamp": 1667941829,
			"wide_cumulative_difficulty": "0x3469a966eb2f788",
			"wide_difficulty": "0x490be69168"
			},
			"credits": 0,
			"json": "{\n  \"major_version\": 16, \n  \"minor_version\": 16, \n  \"timestamp\": 1667941829, \n  \"prev_id\": \"b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7\", \n  \"nonce\": 4110909056, \n  \"miner_tx\": {\n    \"version\": 2, \n    \"unlock_time\": 2751566, \n    \"vin\": [ {\n        \"gen\": {\n          \"height\": 2751506\n        }\n      }\n    ], \n    \"vout\": [ {\n        \"amount\": 600000000000, \n        \"target\": {\n          \"tagged_key\": {\n            \"key\": \"d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85\", \n            \"view_tag\": \"d0\"\n          }\n        }\n      }\n    ], \n    \"extra\": [ 1, 159, 98, 157, 139, 54, 189, 22, 162, 191, 206, 62, 168, 12, 49, 220, 77, 135, 98, 198, 113, 101, 174, 194, 24, 69, 73, 78, 50, 183, 88, 47, 224, 2, 17, 0, 0, 0, 41, 122, 120, 122, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0\n    ], \n    \"rct_signatures\": {\n      \"type\": 0\n    }\n  }, \n  \"tx_hashes\": [ ]\n}",
			"miner_tx_hash": "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
			"status": "OK",
			"top_hash": "",
			"untrusted": false
		}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetBlockResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetBlockResult{
			Blob: "1010c58bab9b06b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7807e07f502cef8a70101ff92f8a7010180e0a596bb1103d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85d034019f629d8b36bd16a2bfce3ea80c31dc4d8762c67165aec21845494e32b7582fe00211000000297a787a000000000000000000000000",
			GetBlockHeaderResult: daemon.GetBlockHeaderResult{
				BlockHeader: daemon.BlockHeader{
					BlockSize:                 106,
					BlockWeight:               106,
					CumulativeDifficulty:      236046001376524168,
					CumulativeDifficultyTop64: 0,
					Depth:                     40,
					Difficulty:                313732272488,
					DifficultyTop64:           0,
					Hash:                      "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428",
					Height:                    2751506,
					LongTermWeight:            176470,
					MajorVersion:              16,
					MinerTxHash:               "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
					MinorVersion:              16,
					Nonce:                     4110909056,
					NumTxes:                   0,
					OrphanStatus:              false,
					PowHash:                   "",
					PrevHash:                  "b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7",
					Reward:                    600000000000,
					Timestamp:                 1667941829,
					WideCumulativeDifficulty:  "0x3469a966eb2f788",
					WideDifficulty:            "0x490be69168"},
				Credits:       0,
				TopHash:       "",
				JsonRpcFooter: defaultMoneroRpcFooter,
			},
			Json: "{\n  \"major_version\": 16, \n  \"minor_version\": 16, \n  \"timestamp\": 1667941829, \n  \"prev_id\": \"b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7\", \n  \"nonce\": 4110909056, \n  \"miner_tx\": {\n    \"version\": 2, \n    \"unlock_time\": 2751566, \n    \"vin\": [ {\n        \"gen\": {\n          \"height\": 2751506\n        }\n      }\n    ], \n    \"vout\": [ {\n        \"amount\": 600000000000, \n        \"target\": {\n          \"tagged_key\": {\n            \"key\": \"d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85\", \n            \"view_tag\": \"d0\"\n          }\n        }\n      }\n    ], \n    \"extra\": [ 1, 159, 98, 157, 139, 54, 189, 22, 162, 191, 206, 62, 168, 12, 49, 220, 77, 135, 98, 198, 113, 101, 174, 194, 24, 69, 73, 78, 50, 183, 88, 47, 224, 2, 17, 0, 0, 0, 41, 122, 120, 122, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0\n    ], \n    \"rct_signatures\": {\n      \"type\": 0\n    }\n  }, \n  \"tx_hashes\": [ ]\n}",
			BlockDetails: daemon.BlockDetails{
				MajorVersion: 16,
				MinorVersion: 16,
				Timestamp:    1667941829,
				PrevId:       "b27bdecfc6cd0a46172d136c08831cf67660377ba992332363228b1b722781e7",
				Nonce:        4110909056,
				MinerTx: daemon.MinerTx{
					Version:    2,
					UnlockTime: 2751566,
					Vin: []daemon.Vin1{
						{
							Gen: daemon.Gen{Height: 2751506},
						},
					},
					Vout: []daemon.Vout1{
						{
							Amount: 600000000000,
							Target: daemon.Target{
								TaggedKey: daemon.TaggedKey{
									Key:     "d7cbf826b665d7a532c316982dc8dbc24f285cbc18bbcc27c7164cd9b3277a85",
									ViewTag: "d0",
								},
							},
						},
					},
					Extra:         []int32{1, 159, 98, 157, 139, 54, 189, 22, 162, 191, 206, 62, 168, 12, 49, 220, 77, 135, 98, 198, 113, 101, 174, 194, 24, 69, 73, 78, 50, 183, 88, 47, 224, 2, 17, 0, 0, 0, 41, 122, 120, 122, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					RctSignatures: daemon.RctSignatures{Type: 0},
				},
				TxHashes: []string{},
			},
			MinerTxHash: "e49b854c5f339d7410a77f2a137281d8042a0ffc7ef9ab24cd670b67139b24cd",
		},
		Error: daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockByHash(exreq.Body.Params.FillPowHash, exreq.Body.Params.Hash)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetTransactionPool(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.EmptyMoneroRpcParams]{Endpoint: "/get_transaction_pool", Body: nil}
	exres := `{
	  "credits": 0,
	  "spent_key_images": [
		{
		  "id_hash": "fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f047873",
		  "txs_hashes": [
			"793da06116f80b9aee790f8558bdfafbc1a7c733ff82f85640d1853dfdc0be4d"
		  ]
		}
	  ],
	  "status": "OK",
	  "top_hash": "",
	  "transactions": [
		{
		  "blob_size": 1529,
		  "do_not_relay": false,
		  "double_spend_seen": false,
		  "fee": 85630000,
		  "id_hash": "793da06116f80b9aee790f8558bdfafbc1a7c733ff82f85640d1853dfdc0be4d",
		  "kept_by_block": false,
		  "last_failed_height": 0,
		  "last_failed_id_hash": "0000000000000000000000000000000000000000000000000000000000000000",
		  "last_relayed_time": 0,
		  "max_used_block_height": 1619940,
		  "max_used_block_id_hash": "f96fd1e6ed43b96250c4a351952de583241a12b86933881f5dc162edd8573a54",
		  "receive_time": 0,
		  "relayed": true,
		  "tx_blob": "020001020010f0c2ca03c5be0af1cb4080d3058db20bdba801d38507f86adc32df2aac04d10aa603f703d00128fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f0478730200030993e6ca2d66871e4869adb2c3a524ad7205fcd3e0b3339daafaea76fc5518ee1b000384f3dd9b4e7df18c5662606a4f6a11ceede3f0cefb41a8586e691baf2930a6fcff2c01b984318d464e56b443af22d5f880470606435172a3bad71966e2a4bae5d18a8002090190d13c4c7d9222d206b0b8ea288b2fe303da838a84779fe795bc0ba77509cd23fa0e8ec03a348ed6e80386c93c276ef69f1c223f811ffc6ce1e88c030a28ceaa373700ce1aeb4167861ec41494edf53f3d7b7568fa7ac05db0aaf324da012644a5380b8de1652a3d47654ecee118eca9506655e77fb0e339aef31da452dd360227a720fca490111110bb23126a49cf783cb67ab8cd91de4891db2e7898ab6923bc04f5917dbe17dd5e6ef9d248cd7bb01afb4675eef4bc8fb7707c7a470ae1bd93860a4ad45f2d1ca2bddefa4f1598cf20be56051cae5b61c3f379f6160e1298b6aeaa25fdfb8631a32dd9bf8efeb66387304516e8bd00599caaa8a77104600b39b3e3f9390e7f6cb61062021d2e8d7f6a6fbb7b04318f35077a3243390f07b8bdee2c2c2997f46c9dd024f5bad1004a52cc8cbcb051fcc63de46476962fd79bce03d001b6ed12d6417ae5871e2a05574316ac53050712cf4129c5b00534798facd82baf29aaa8a96dd3e04cf6c742544b3aa6b37bce394c416869be0bb145f64be9871eda186cdce9c8fefec3cda6a70574492c42ff4c998e82f494192f02f98a7ecc762c59608409508924bed2665b53c20b93fb3338c2edad582ca19ef77cc02f17f547386b014b1ad6a79df59130f71c05cf7f50abd447c01249afdd7ffafdf6f43138b4905838243884fe16216df87300e1bf5e20e78ecea69bc53e1a07c2da698b34dce738ec74a2cba0b130378d1cf15a3697566a59bbcea9a082cd16e72907754e50b6b3daa866f459634f8e53ba531953c227309cf8f7a7fbfaac2daa5a4811b347f89eb981f331b752313aa8dc7aad366a40bc3e2ef68c51733e0e228769927c8d8eaccd0640a02916604234e7a1b1cc7f7e8311815452668becfc3d76332ea1de6ee160660fc310148d49135b718e611d1ade4146dd813253928721c48f76ca5d59d19b257afdd8c2d5abfe1c905ec00c34d150b90a52683c58d33506f70f64346d5ca69a26007689eb79755e9953f21bce011087d065ca137e4bdfae579e248336f3d39f4a880823b68e571ca8c3adbbd91fb90be2f5c7832007b39e788f94f3ccd48dc6b09d87d3b3d71c0a6df53658969b5a18d7864be6a00ab356d93b50cd3aae005c891cb72047726b7a40228bd1ac547f08b0ba2b5b630a693582bb3a5e39ebe2a66b44d5fe856875efffec516e2ca5229fb9689a92c1087cfabb788fc5925f23a45b675e28ff696009d928d25e3edce01703135ffc6404159297800e32b019ee70b15e73d4d91d4c439ad13bde42eee8f59120aedf0607b95ba55a6497a52e476718d0f4c8353190418fe6b2f4cc7050ced06451fb6d049e92a46ad7d55fe6aaf07faa17d791d7ee8ca2ac49e98417392575857bcdc206c72d57a1933434c5cd8b5fa167cb7d8b512347956fd6bc60caeb269f30beb60e5991d37f9543d81b0cd4a04087b8fbc19eb98102d3b460608da705354ac28a0a923382f6792d746b9c7bc5f7f00b01bebcf3a173c78c268872feb49422d8840e541f7c83b4da45bf3289eb36772444e08e703347313ab0500614c8b571b35d07279006100ed62a32e592071e8e749895090e27c347f2567bfbace5a7823100007b29c0c7d11657ead227902d6a95e855cf38a63bdd963fe99f80c7a5da27fc0b7f7fd35f789b110cac086707a498f03b692ec210a2a52f90114827bb8b53da058f443440db05a72ccaa68ac8cc022b067e122c563b5c277703fecac7bb876609ef5c502d5ab8701c613b7ee3ed20069681e0e98b54169e4a0e2f165ee1fc9e0e6213c0f6e752de084e9f90a492d1a5b42fe2b82ebd4f1f1228dfcaea591d4271c39fb4de4c13906eb11eb2da196165ac075f5d797301cb5f88e80023532a063f",
		  "tx_json": "{\n  \"version\": 2, \n  \"unlock_time\": 0, \n  \"vin\": [ {\n      \"key\": {\n        \"amount\": 0, \n        \"key_offsets\": [ 7512432, 171845, 1058289, 92544, 186637, 21595, 115411, 13688, 6492, 5471, 556, 1361, 422, 503, 208, 40\n        ], \n        \"k_image\": \"fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f047873\"\n      }\n    }\n  ], \n  \"vout\": [ {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"0993e6ca2d66871e4869adb2c3a524ad7205fcd3e0b3339daafaea76fc5518ee\", \n          \"view_tag\": \"1b\"\n        }\n      }\n    }, {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"84f3dd9b4e7df18c5662606a4f6a11ceede3f0cefb41a8586e691baf2930a6fc\", \n          \"view_tag\": \"ff\"\n        }\n      }\n    }\n  ], \n  \"extra\": [ 1, 185, 132, 49, 141, 70, 78, 86, 180, 67, 175, 34, 213, 248, 128, 71, 6, 6, 67, 81, 114, 163, 186, 215, 25, 102, 226, 164, 186, 229, 209, 138, 128, 2, 9, 1, 144, 209, 60, 76, 125, 146, 34, 210\n  ], \n  \"rct_signatures\": {\n    \"type\": 6, \n    \"txnFee\": 85630000, \n    \"ecdhInfo\": [ {\n        \"amount\": \"8b2fe303da838a84\"\n      }, {\n        \"amount\": \"779fe795bc0ba775\"\n      }], \n    \"outPk\": [ \"09cd23fa0e8ec03a348ed6e80386c93c276ef69f1c223f811ffc6ce1e88c030a\", \"28ceaa373700ce1aeb4167861ec41494edf53f3d7b7568fa7ac05db0aaf324da\"]\n  }, \n  \"rctsig_prunable\": {\n    \"nbp\": 1, \n    \"bpp\": [ {\n        \"A\": \"2644a5380b8de1652a3d47654ecee118eca9506655e77fb0e339aef31da452dd\", \n        \"A1\": \"360227a720fca490111110bb23126a49cf783cb67ab8cd91de4891db2e7898ab\", \n        \"B\": \"6923bc04f5917dbe17dd5e6ef9d248cd7bb01afb4675eef4bc8fb7707c7a470a\", \n        \"r1\": \"e1bd93860a4ad45f2d1ca2bddefa4f1598cf20be56051cae5b61c3f379f6160e\", \n        \"s1\": \"1298b6aeaa25fdfb8631a32dd9bf8efeb66387304516e8bd00599caaa8a77104\", \n        \"d1\": \"600b39b3e3f9390e7f6cb61062021d2e8d7f6a6fbb7b04318f35077a3243390f\", \n        \"L\": [ \"b8bdee2c2c2997f46c9dd024f5bad1004a52cc8cbcb051fcc63de46476962fd7\", \"9bce03d001b6ed12d6417ae5871e2a05574316ac53050712cf4129c5b0053479\", \"8facd82baf29aaa8a96dd3e04cf6c742544b3aa6b37bce394c416869be0bb145\", \"f64be9871eda186cdce9c8fefec3cda6a70574492c42ff4c998e82f494192f02\", \"f98a7ecc762c59608409508924bed2665b53c20b93fb3338c2edad582ca19ef7\", \"7cc02f17f547386b014b1ad6a79df59130f71c05cf7f50abd447c01249afdd7f\", \"fafdf6f43138b4905838243884fe16216df87300e1bf5e20e78ecea69bc53e1a\"\n        ], \n        \"R\": [ \"c2da698b34dce738ec74a2cba0b130378d1cf15a3697566a59bbcea9a082cd16\", \"e72907754e50b6b3daa866f459634f8e53ba531953c227309cf8f7a7fbfaac2d\", \"aa5a4811b347f89eb981f331b752313aa8dc7aad366a40bc3e2ef68c51733e0e\", \"228769927c8d8eaccd0640a02916604234e7a1b1cc7f7e8311815452668becfc\", \"3d76332ea1de6ee160660fc310148d49135b718e611d1ade4146dd8132539287\", \"21c48f76ca5d59d19b257afdd8c2d5abfe1c905ec00c34d150b90a52683c58d3\", \"3506f70f64346d5ca69a26007689eb79755e9953f21bce011087d065ca137e4b\"\n        ]\n      }\n    ], \n    \"CLSAGs\": [ {\n        \"s\": [ \"dfae579e248336f3d39f4a880823b68e571ca8c3adbbd91fb90be2f5c7832007\", \"b39e788f94f3ccd48dc6b09d87d3b3d71c0a6df53658969b5a18d7864be6a00a\", \"b356d93b50cd3aae005c891cb72047726b7a40228bd1ac547f08b0ba2b5b630a\", \"693582bb3a5e39ebe2a66b44d5fe856875efffec516e2ca5229fb9689a92c108\", \"7cfabb788fc5925f23a45b675e28ff696009d928d25e3edce01703135ffc6404\", \"159297800e32b019ee70b15e73d4d91d4c439ad13bde42eee8f59120aedf0607\", \"b95ba55a6497a52e476718d0f4c8353190418fe6b2f4cc7050ced06451fb6d04\", \"9e92a46ad7d55fe6aaf07faa17d791d7ee8ca2ac49e98417392575857bcdc206\", \"c72d57a1933434c5cd8b5fa167cb7d8b512347956fd6bc60caeb269f30beb60e\", \"5991d37f9543d81b0cd4a04087b8fbc19eb98102d3b460608da705354ac28a0a\", \"923382f6792d746b9c7bc5f7f00b01bebcf3a173c78c268872feb49422d8840e\", \"541f7c83b4da45bf3289eb36772444e08e703347313ab0500614c8b571b35d07\", \"279006100ed62a32e592071e8e749895090e27c347f2567bfbace5a782310000\", \"7b29c0c7d11657ead227902d6a95e855cf38a63bdd963fe99f80c7a5da27fc0b\", \"7f7fd35f789b110cac086707a498f03b692ec210a2a52f90114827bb8b53da05\", \"8f443440db05a72ccaa68ac8cc022b067e122c563b5c277703fecac7bb876609\"], \n        \"c1\": \"ef5c502d5ab8701c613b7ee3ed20069681e0e98b54169e4a0e2f165ee1fc9e0e\", \n        \"D\": \"6213c0f6e752de084e9f90a492d1a5b42fe2b82ebd4f1f1228dfcaea591d4271\"\n      }], \n    \"pseudoOuts\": [ \"c39fb4de4c13906eb11eb2da196165ac075f5d797301cb5f88e80023532a063f\"]\n  }\n}",
		  "weight": 1529
		}
	  ],
	  "untrusted": false
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.GetTransactionPoolResponse{
		Credits: 0,
		SpentKeyImages: []daemon.SpentKeyImage{
			{
				IdHash:    "fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f047873",
				TxsHashes: []string{"793da06116f80b9aee790f8558bdfafbc1a7c733ff82f85640d1853dfdc0be4d"},
			},
		},
		TopHash: "",
		Transactions: []daemon.MoneroTx{
			{
				BlobSize:           1529,
				DoNotRelay:         false,
				DoubleSpendSeen:    false,
				Fee:                85630000,
				IdHash:             "793da06116f80b9aee790f8558bdfafbc1a7c733ff82f85640d1853dfdc0be4d",
				KeptByBlock:        false,
				LastFailedHeight:   0,
				LastFailedIdHash:   "0000000000000000000000000000000000000000000000000000000000000000",
				LastRelayedTime:    0,
				MaxUsedBlockHeight: 1619940,
				MaxUsedBlockIdHash: "f96fd1e6ed43b96250c4a351952de583241a12b86933881f5dc162edd8573a54",
				ReceiveTime:        0,
				Relayed:            true,
				TxBlob:             "020001020010f0c2ca03c5be0af1cb4080d3058db20bdba801d38507f86adc32df2aac04d10aa603f703d00128fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f0478730200030993e6ca2d66871e4869adb2c3a524ad7205fcd3e0b3339daafaea76fc5518ee1b000384f3dd9b4e7df18c5662606a4f6a11ceede3f0cefb41a8586e691baf2930a6fcff2c01b984318d464e56b443af22d5f880470606435172a3bad71966e2a4bae5d18a8002090190d13c4c7d9222d206b0b8ea288b2fe303da838a84779fe795bc0ba77509cd23fa0e8ec03a348ed6e80386c93c276ef69f1c223f811ffc6ce1e88c030a28ceaa373700ce1aeb4167861ec41494edf53f3d7b7568fa7ac05db0aaf324da012644a5380b8de1652a3d47654ecee118eca9506655e77fb0e339aef31da452dd360227a720fca490111110bb23126a49cf783cb67ab8cd91de4891db2e7898ab6923bc04f5917dbe17dd5e6ef9d248cd7bb01afb4675eef4bc8fb7707c7a470ae1bd93860a4ad45f2d1ca2bddefa4f1598cf20be56051cae5b61c3f379f6160e1298b6aeaa25fdfb8631a32dd9bf8efeb66387304516e8bd00599caaa8a77104600b39b3e3f9390e7f6cb61062021d2e8d7f6a6fbb7b04318f35077a3243390f07b8bdee2c2c2997f46c9dd024f5bad1004a52cc8cbcb051fcc63de46476962fd79bce03d001b6ed12d6417ae5871e2a05574316ac53050712cf4129c5b00534798facd82baf29aaa8a96dd3e04cf6c742544b3aa6b37bce394c416869be0bb145f64be9871eda186cdce9c8fefec3cda6a70574492c42ff4c998e82f494192f02f98a7ecc762c59608409508924bed2665b53c20b93fb3338c2edad582ca19ef77cc02f17f547386b014b1ad6a79df59130f71c05cf7f50abd447c01249afdd7ffafdf6f43138b4905838243884fe16216df87300e1bf5e20e78ecea69bc53e1a07c2da698b34dce738ec74a2cba0b130378d1cf15a3697566a59bbcea9a082cd16e72907754e50b6b3daa866f459634f8e53ba531953c227309cf8f7a7fbfaac2daa5a4811b347f89eb981f331b752313aa8dc7aad366a40bc3e2ef68c51733e0e228769927c8d8eaccd0640a02916604234e7a1b1cc7f7e8311815452668becfc3d76332ea1de6ee160660fc310148d49135b718e611d1ade4146dd813253928721c48f76ca5d59d19b257afdd8c2d5abfe1c905ec00c34d150b90a52683c58d33506f70f64346d5ca69a26007689eb79755e9953f21bce011087d065ca137e4bdfae579e248336f3d39f4a880823b68e571ca8c3adbbd91fb90be2f5c7832007b39e788f94f3ccd48dc6b09d87d3b3d71c0a6df53658969b5a18d7864be6a00ab356d93b50cd3aae005c891cb72047726b7a40228bd1ac547f08b0ba2b5b630a693582bb3a5e39ebe2a66b44d5fe856875efffec516e2ca5229fb9689a92c1087cfabb788fc5925f23a45b675e28ff696009d928d25e3edce01703135ffc6404159297800e32b019ee70b15e73d4d91d4c439ad13bde42eee8f59120aedf0607b95ba55a6497a52e476718d0f4c8353190418fe6b2f4cc7050ced06451fb6d049e92a46ad7d55fe6aaf07faa17d791d7ee8ca2ac49e98417392575857bcdc206c72d57a1933434c5cd8b5fa167cb7d8b512347956fd6bc60caeb269f30beb60e5991d37f9543d81b0cd4a04087b8fbc19eb98102d3b460608da705354ac28a0a923382f6792d746b9c7bc5f7f00b01bebcf3a173c78c268872feb49422d8840e541f7c83b4da45bf3289eb36772444e08e703347313ab0500614c8b571b35d07279006100ed62a32e592071e8e749895090e27c347f2567bfbace5a7823100007b29c0c7d11657ead227902d6a95e855cf38a63bdd963fe99f80c7a5da27fc0b7f7fd35f789b110cac086707a498f03b692ec210a2a52f90114827bb8b53da058f443440db05a72ccaa68ac8cc022b067e122c563b5c277703fecac7bb876609ef5c502d5ab8701c613b7ee3ed20069681e0e98b54169e4a0e2f165ee1fc9e0e6213c0f6e752de084e9f90a492d1a5b42fe2b82ebd4f1f1228dfcaea591d4271c39fb4de4c13906eb11eb2da196165ac075f5d797301cb5f88e80023532a063f",
				TxJson:             "{\n  \"version\": 2, \n  \"unlock_time\": 0, \n  \"vin\": [ {\n      \"key\": {\n        \"amount\": 0, \n        \"key_offsets\": [ 7512432, 171845, 1058289, 92544, 186637, 21595, 115411, 13688, 6492, 5471, 556, 1361, 422, 503, 208, 40\n        ], \n        \"k_image\": \"fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f047873\"\n      }\n    }\n  ], \n  \"vout\": [ {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"0993e6ca2d66871e4869adb2c3a524ad7205fcd3e0b3339daafaea76fc5518ee\", \n          \"view_tag\": \"1b\"\n        }\n      }\n    }, {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"84f3dd9b4e7df18c5662606a4f6a11ceede3f0cefb41a8586e691baf2930a6fc\", \n          \"view_tag\": \"ff\"\n        }\n      }\n    }\n  ], \n  \"extra\": [ 1, 185, 132, 49, 141, 70, 78, 86, 180, 67, 175, 34, 213, 248, 128, 71, 6, 6, 67, 81, 114, 163, 186, 215, 25, 102, 226, 164, 186, 229, 209, 138, 128, 2, 9, 1, 144, 209, 60, 76, 125, 146, 34, 210\n  ], \n  \"rct_signatures\": {\n    \"type\": 6, \n    \"txnFee\": 85630000, \n    \"ecdhInfo\": [ {\n        \"amount\": \"8b2fe303da838a84\"\n      }, {\n        \"amount\": \"779fe795bc0ba775\"\n      }], \n    \"outPk\": [ \"09cd23fa0e8ec03a348ed6e80386c93c276ef69f1c223f811ffc6ce1e88c030a\", \"28ceaa373700ce1aeb4167861ec41494edf53f3d7b7568fa7ac05db0aaf324da\"]\n  }, \n  \"rctsig_prunable\": {\n    \"nbp\": 1, \n    \"bpp\": [ {\n        \"A\": \"2644a5380b8de1652a3d47654ecee118eca9506655e77fb0e339aef31da452dd\", \n        \"A1\": \"360227a720fca490111110bb23126a49cf783cb67ab8cd91de4891db2e7898ab\", \n        \"B\": \"6923bc04f5917dbe17dd5e6ef9d248cd7bb01afb4675eef4bc8fb7707c7a470a\", \n        \"r1\": \"e1bd93860a4ad45f2d1ca2bddefa4f1598cf20be56051cae5b61c3f379f6160e\", \n        \"s1\": \"1298b6aeaa25fdfb8631a32dd9bf8efeb66387304516e8bd00599caaa8a77104\", \n        \"d1\": \"600b39b3e3f9390e7f6cb61062021d2e8d7f6a6fbb7b04318f35077a3243390f\", \n        \"L\": [ \"b8bdee2c2c2997f46c9dd024f5bad1004a52cc8cbcb051fcc63de46476962fd7\", \"9bce03d001b6ed12d6417ae5871e2a05574316ac53050712cf4129c5b0053479\", \"8facd82baf29aaa8a96dd3e04cf6c742544b3aa6b37bce394c416869be0bb145\", \"f64be9871eda186cdce9c8fefec3cda6a70574492c42ff4c998e82f494192f02\", \"f98a7ecc762c59608409508924bed2665b53c20b93fb3338c2edad582ca19ef7\", \"7cc02f17f547386b014b1ad6a79df59130f71c05cf7f50abd447c01249afdd7f\", \"fafdf6f43138b4905838243884fe16216df87300e1bf5e20e78ecea69bc53e1a\"\n        ], \n        \"R\": [ \"c2da698b34dce738ec74a2cba0b130378d1cf15a3697566a59bbcea9a082cd16\", \"e72907754e50b6b3daa866f459634f8e53ba531953c227309cf8f7a7fbfaac2d\", \"aa5a4811b347f89eb981f331b752313aa8dc7aad366a40bc3e2ef68c51733e0e\", \"228769927c8d8eaccd0640a02916604234e7a1b1cc7f7e8311815452668becfc\", \"3d76332ea1de6ee160660fc310148d49135b718e611d1ade4146dd8132539287\", \"21c48f76ca5d59d19b257afdd8c2d5abfe1c905ec00c34d150b90a52683c58d3\", \"3506f70f64346d5ca69a26007689eb79755e9953f21bce011087d065ca137e4b\"\n        ]\n      }\n    ], \n    \"CLSAGs\": [ {\n        \"s\": [ \"dfae579e248336f3d39f4a880823b68e571ca8c3adbbd91fb90be2f5c7832007\", \"b39e788f94f3ccd48dc6b09d87d3b3d71c0a6df53658969b5a18d7864be6a00a\", \"b356d93b50cd3aae005c891cb72047726b7a40228bd1ac547f08b0ba2b5b630a\", \"693582bb3a5e39ebe2a66b44d5fe856875efffec516e2ca5229fb9689a92c108\", \"7cfabb788fc5925f23a45b675e28ff696009d928d25e3edce01703135ffc6404\", \"159297800e32b019ee70b15e73d4d91d4c439ad13bde42eee8f59120aedf0607\", \"b95ba55a6497a52e476718d0f4c8353190418fe6b2f4cc7050ced06451fb6d04\", \"9e92a46ad7d55fe6aaf07faa17d791d7ee8ca2ac49e98417392575857bcdc206\", \"c72d57a1933434c5cd8b5fa167cb7d8b512347956fd6bc60caeb269f30beb60e\", \"5991d37f9543d81b0cd4a04087b8fbc19eb98102d3b460608da705354ac28a0a\", \"923382f6792d746b9c7bc5f7f00b01bebcf3a173c78c268872feb49422d8840e\", \"541f7c83b4da45bf3289eb36772444e08e703347313ab0500614c8b571b35d07\", \"279006100ed62a32e592071e8e749895090e27c347f2567bfbace5a782310000\", \"7b29c0c7d11657ead227902d6a95e855cf38a63bdd963fe99f80c7a5da27fc0b\", \"7f7fd35f789b110cac086707a498f03b692ec210a2a52f90114827bb8b53da05\", \"8f443440db05a72ccaa68ac8cc022b067e122c563b5c277703fecac7bb876609\"], \n        \"c1\": \"ef5c502d5ab8701c613b7ee3ed20069681e0e98b54169e4a0e2f165ee1fc9e0e\", \n        \"D\": \"6213c0f6e752de084e9f90a492d1a5b42fe2b82ebd4f1f1228dfcaea591d4271\"\n      }], \n    \"pseudoOuts\": [ \"c39fb4de4c13906eb11eb2da196165ac075f5d797301cb5f88e80023532a063f\"]\n  }\n}",
				TxInfo: daemon.MoneroTxInfo{
					Version:    2,
					UnlockTime: 0,
					Vin: []daemon.Vin2{
						{
							Key: daemon.Key{
								Amount:     0,
								KeyOffsets: []int64{7512432, 171845, 1058289, 92544, 186637, 21595, 115411, 13688, 6492, 5471, 556, 1361, 422, 503, 208, 40},
								KeyImage:   "fc5655d843ed8b30a3563bbff1d02b606b089b1725c717823b0898c52f047873",
							},
						},
					},
					Vout: []daemon.Vout1{
						{
							Amount: 0,
							Target: daemon.Target{
								TaggedKey: daemon.TaggedKey{
									Key:     "0993e6ca2d66871e4869adb2c3a524ad7205fcd3e0b3339daafaea76fc5518ee",
									ViewTag: "1b",
								},
							},
						},
						{
							Amount: 0,
							Target: daemon.Target{
								TaggedKey: daemon.TaggedKey{
									Key:     "84f3dd9b4e7df18c5662606a4f6a11ceede3f0cefb41a8586e691baf2930a6fc",
									ViewTag: "ff",
								},
							},
						},
					},
					Extra: []int32{1, 185, 132, 49, 141, 70, 78, 86, 180, 67, 175, 34, 213, 248, 128, 71, 6, 6, 67, 81, 114, 163, 186, 215, 25, 102, 226, 164, 186, 229, 209, 138, 128, 2, 9, 1, 144, 209, 60, 76, 125, 146, 34, 210},
					RctSignatures: daemon.RctSignature{
						Type:   6,
						TxnFee: 85630000,
						EcdhInfo: []daemon.EcdhInfo{
							{
								Amount: "8b2fe303da838a84",
							},
							{
								Amount: "779fe795bc0ba775",
							},
						},
						OutPk: []string{"09cd23fa0e8ec03a348ed6e80386c93c276ef69f1c223f811ffc6ce1e88c030a", "28ceaa373700ce1aeb4167861ec41494edf53f3d7b7568fa7ac05db0aaf324da"},
					},
					RctsigPrunable: daemon.RctsigPrunable{
						Nbp: 1,
						Bpp: []daemon.Bpp{
							{
								A:  "2644a5380b8de1652a3d47654ecee118eca9506655e77fb0e339aef31da452dd",
								A1: "360227a720fca490111110bb23126a49cf783cb67ab8cd91de4891db2e7898ab",
								B:  "6923bc04f5917dbe17dd5e6ef9d248cd7bb01afb4675eef4bc8fb7707c7a470a",
								R1: "e1bd93860a4ad45f2d1ca2bddefa4f1598cf20be56051cae5b61c3f379f6160e",
								S1: "1298b6aeaa25fdfb8631a32dd9bf8efeb66387304516e8bd00599caaa8a77104",
								D1: "600b39b3e3f9390e7f6cb61062021d2e8d7f6a6fbb7b04318f35077a3243390f",
								L:  []string{"b8bdee2c2c2997f46c9dd024f5bad1004a52cc8cbcb051fcc63de46476962fd7", "9bce03d001b6ed12d6417ae5871e2a05574316ac53050712cf4129c5b0053479", "8facd82baf29aaa8a96dd3e04cf6c742544b3aa6b37bce394c416869be0bb145", "f64be9871eda186cdce9c8fefec3cda6a70574492c42ff4c998e82f494192f02", "f98a7ecc762c59608409508924bed2665b53c20b93fb3338c2edad582ca19ef7", "7cc02f17f547386b014b1ad6a79df59130f71c05cf7f50abd447c01249afdd7f", "fafdf6f43138b4905838243884fe16216df87300e1bf5e20e78ecea69bc53e1a"},
								R:  []string{"c2da698b34dce738ec74a2cba0b130378d1cf15a3697566a59bbcea9a082cd16", "e72907754e50b6b3daa866f459634f8e53ba531953c227309cf8f7a7fbfaac2d", "aa5a4811b347f89eb981f331b752313aa8dc7aad366a40bc3e2ef68c51733e0e", "228769927c8d8eaccd0640a02916604234e7a1b1cc7f7e8311815452668becfc", "3d76332ea1de6ee160660fc310148d49135b718e611d1ade4146dd8132539287", "21c48f76ca5d59d19b257afdd8c2d5abfe1c905ec00c34d150b90a52683c58d3", "3506f70f64346d5ca69a26007689eb79755e9953f21bce011087d065ca137e4b"},
							},
						},
						CLSAGs: []daemon.CLSAG{
							{
								S:  []string{"dfae579e248336f3d39f4a880823b68e571ca8c3adbbd91fb90be2f5c7832007", "b39e788f94f3ccd48dc6b09d87d3b3d71c0a6df53658969b5a18d7864be6a00a", "b356d93b50cd3aae005c891cb72047726b7a40228bd1ac547f08b0ba2b5b630a", "693582bb3a5e39ebe2a66b44d5fe856875efffec516e2ca5229fb9689a92c108", "7cfabb788fc5925f23a45b675e28ff696009d928d25e3edce01703135ffc6404", "159297800e32b019ee70b15e73d4d91d4c439ad13bde42eee8f59120aedf0607", "b95ba55a6497a52e476718d0f4c8353190418fe6b2f4cc7050ced06451fb6d04", "9e92a46ad7d55fe6aaf07faa17d791d7ee8ca2ac49e98417392575857bcdc206", "c72d57a1933434c5cd8b5fa167cb7d8b512347956fd6bc60caeb269f30beb60e", "5991d37f9543d81b0cd4a04087b8fbc19eb98102d3b460608da705354ac28a0a", "923382f6792d746b9c7bc5f7f00b01bebcf3a173c78c268872feb49422d8840e", "541f7c83b4da45bf3289eb36772444e08e703347313ab0500614c8b571b35d07", "279006100ed62a32e592071e8e749895090e27c347f2567bfbace5a782310000", "7b29c0c7d11657ead227902d6a95e855cf38a63bdd963fe99f80c7a5da27fc0b", "7f7fd35f789b110cac086707a498f03b692ec210a2a52f90114827bb8b53da05", "8f443440db05a72ccaa68ac8cc022b067e122c563b5c277703fecac7bb876609"},
								C1: "ef5c502d5ab8701c613b7ee3ed20069681e0e98b54169e4a0e2f165ee1fc9e0e",
								D:  "6213c0f6e752de084e9f90a492d1a5b42fe2b82ebd4f1f1228dfcaea591d4271",
							},
						},
						PseudoOuts: []string{"c39fb4de4c13906eb11eb2da196165ac075f5d797301cb5f88e80023532a063f"},
					},
				},
				Weight: 1529,
			},
		},
		JsonRpcFooter: defaultMoneroRpcFooter,
	}

	actual, err := test_daemon.GetTransactionPool()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expected, actual)
}

func TestGetTransactions(t *testing.T) {
	reqBody := daemon.GetTransactionsParams{
		TxHashes:     []string{"d6e48158472848e6687173a91ae6eebfa3e1d778e65252ee99d7515d63090408"},
		DecodeAsJson: true,
		Prune:        false,
		Split:        false,
	}
	exreq := &daemon.MoneroRpcRequest[daemon.GetTransactionsParams]{Endpoint: "/get_transactions", Body: &reqBody}
	exres := `{
		"credits": 0,
		"status": "OK",
		"top_hash": "",
		"txs": [
			{
				"as_hex": "020001020010899cde23f4c8800784cf8e07e2af04d48a0cbdb142ccbb0c965bf79901bd8505bbfb01a18401a6a60189138034cf11045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5020003eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65abbd00032dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef4c2c011dead87c2a407674d80d93b7804384ad59361d3b01629dce957ebe23fdbd5205020901c43add31d092c19b0680f1d03a401150d04a7559e8807036489f864b379f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a31631951589265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d470601f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a57724e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c073aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd4796845d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a654641793655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8db2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba590707849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c94018e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e306706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c78490316e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc3440cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb22038094443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab037e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0bbef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170ed60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af005011960059614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810dc8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466",
				"as_json": "{\n  \"version\": 2, \n  \"unlock_time\": 0, \n  \"vin\": [ {\n      \"key\": {\n        \"amount\": 0, \n        \"key_offsets\": [ 74944009, 14689396, 14919556, 71650, 197972, 1087677, 204236, 11670, 19703, 82621, 32187, 16929, 21286, 2441, 6656, 2255\n        ], \n        \"k_image\": \"045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5\"\n      }\n    }\n  ], \n  \"vout\": [ {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65ab\", \n          \"view_tag\": \"bd\"\n        }\n      }\n    }, {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"2dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef\", \n          \"view_tag\": \"4c\"\n        }\n      }\n    }\n  ], \n  \"extra\": [ 1, 29, 234, 216, 124, 42, 64, 118, 116, 216, 13, 147, 183, 128, 67, 132, 173, 89, 54, 29, 59, 1, 98, 157, 206, 149, 126, 190, 35, 253, 189, 82, 5, 2, 9, 1, 196, 58, 221, 49, 208, 146, 193, 155\n  ], \n  \"rct_signatures\": {\n    \"type\": 6, \n    \"txnFee\": 122960000, \n    \"ecdhInfo\": [ {\n        \"trunc_amount\": \"401150d04a7559e8\"\n      }, {\n        \"trunc_amount\": \"807036489f864b37\"\n      }], \n    \"outPk\": [ \"9f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a3163195158\", \"9265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d4706\"]\n  }, \n  \"rctsig_prunable\": {\n    \"nbp\": 1, \n    \"bpp\": [ {\n        \"A\": \"f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a5772\", \n        \"A1\": \"4e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b\", \n        \"B\": \"68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc\", \n        \"r1\": \"5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100\", \n        \"s1\": \"bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00\", \n        \"d1\": \"cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c\", \n        \"L\": [ \"3aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d\", \"881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e\", \"2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd47968\", \"45d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a65464179\", \"3655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8d\", \"b2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13\", \"c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba59\"\n        ], \n        \"R\": [ \"07849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c9401\", \"8e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e30\", \"6706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d\", \"880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c7849031\", \"6e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117\", \"310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc34\", \"40cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a\"\n        ]\n      }\n    ], \n    \"CLSAGs\": [ {\n        \"s\": [ \"0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb2203809\", \"4443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a\", \"40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e\", \"7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab03\", \"7e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03\", \"d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303\", \"f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402\", \"ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0b\", \"bef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07\", \"e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170e\", \"d60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401\", \"b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d\", \"2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507\", \"920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af00501196005\", \"9614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07\", \"347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810d\"], \n        \"c1\": \"c8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205\", \n        \"D\": \"b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a\"\n      }], \n    \"pseudoOuts\": [ \"6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466\"]\n  }\n}",
				"block_height": 3169795,
				"block_timestamp": 1718210909,
				"confirmations": 24,
				"double_spend_seen": false,
				"in_pool": false,
				"output_indices": [
					106312029,
					106312030
				],
				"prunable_as_hex": "",
				"prunable_hash": "734646af580b1e38eea27214d821670c731f92be1bea142bbb5a90b923a18292",
				"pruned_as_hex": "",
				"tx_hash": "45b27c7cde61cdbb05b957241a4cc698665f7fbc55af62c84cab45f71b879a15"
			}
		],
		"txs_as_hex": [
			"020001020010899cde23f4c8800784cf8e07e2af04d48a0cbdb142ccbb0c965bf79901bd8505bbfb01a18401a6a60189138034cf11045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5020003eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65abbd00032dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef4c2c011dead87c2a407674d80d93b7804384ad59361d3b01629dce957ebe23fdbd5205020901c43add31d092c19b0680f1d03a401150d04a7559e8807036489f864b379f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a31631951589265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d470601f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a57724e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c073aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd4796845d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a654641793655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8db2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba590707849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c94018e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e306706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c78490316e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc3440cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb22038094443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab037e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0bbef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170ed60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af005011960059614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810dc8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466"
		],
		"txs_as_json": [
			"{\n  \"version\": 2, \n  \"unlock_time\": 0, \n  \"vin\": [ {\n      \"key\": {\n        \"amount\": 0, \n        \"key_offsets\": [ 74944009, 14689396, 14919556, 71650, 197972, 1087677, 204236, 11670, 19703, 82621, 32187, 16929, 21286, 2441, 6656, 2255\n        ], \n        \"k_image\": \"045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5\"\n      }\n    }\n  ], \n  \"vout\": [ {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65ab\", \n          \"view_tag\": \"bd\"\n        }\n      }\n    }, {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"2dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef\", \n          \"view_tag\": \"4c\"\n        }\n      }\n    }\n  ], \n  \"extra\": [ 1, 29, 234, 216, 124, 42, 64, 118, 116, 216, 13, 147, 183, 128, 67, 132, 173, 89, 54, 29, 59, 1, 98, 157, 206, 149, 126, 190, 35, 253, 189, 82, 5, 2, 9, 1, 196, 58, 221, 49, 208, 146, 193, 155\n  ], \n  \"rct_signatures\": {\n    \"type\": 6, \n    \"txnFee\": 122960000, \n    \"ecdhInfo\": [ {\n        \"trunc_amount\": \"401150d04a7559e8\"\n      }, {\n        \"trunc_amount\": \"807036489f864b37\"\n      }], \n    \"outPk\": [ \"9f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a3163195158\", \"9265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d4706\"]\n  }, \n  \"rctsig_prunable\": {\n    \"nbp\": 1, \n    \"bpp\": [ {\n        \"A\": \"f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a5772\", \n        \"A1\": \"4e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b\", \n        \"B\": \"68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc\", \n        \"r1\": \"5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100\", \n        \"s1\": \"bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00\", \n        \"d1\": \"cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c\", \n        \"L\": [ \"3aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d\", \"881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e\", \"2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd47968\", \"45d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a65464179\", \"3655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8d\", \"b2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13\", \"c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba59\"\n        ], \n        \"R\": [ \"07849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c9401\", \"8e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e30\", \"6706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d\", \"880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c7849031\", \"6e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117\", \"310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc34\", \"40cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a\"\n        ]\n      }\n    ], \n    \"CLSAGs\": [ {\n        \"s\": [ \"0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb2203809\", \"4443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a\", \"40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e\", \"7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab03\", \"7e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03\", \"d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303\", \"f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402\", \"ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0b\", \"bef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07\", \"e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170e\", \"d60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401\", \"b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d\", \"2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507\", \"920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af00501196005\", \"9614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07\", \"347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810d\"], \n        \"c1\": \"c8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205\", \n        \"D\": \"b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a\"\n      }], \n    \"pseudoOuts\": [ \"6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466\"]\n  }\n}"
		],
		"untrusted": false
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.GetTransactionsResponse{
		Credits:  0,
		MissedTx: nil,
		TopHash:  "",
		Txs: []daemon.MoneroTx1{
			{
				AsHex:           "020001020010899cde23f4c8800784cf8e07e2af04d48a0cbdb142ccbb0c965bf79901bd8505bbfb01a18401a6a60189138034cf11045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5020003eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65abbd00032dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef4c2c011dead87c2a407674d80d93b7804384ad59361d3b01629dce957ebe23fdbd5205020901c43add31d092c19b0680f1d03a401150d04a7559e8807036489f864b379f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a31631951589265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d470601f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a57724e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c073aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd4796845d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a654641793655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8db2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba590707849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c94018e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e306706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c78490316e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc3440cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb22038094443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab037e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0bbef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170ed60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af005011960059614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810dc8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466",
				AsJson:          "{\n  \"version\": 2, \n  \"unlock_time\": 0, \n  \"vin\": [ {\n      \"key\": {\n        \"amount\": 0, \n        \"key_offsets\": [ 74944009, 14689396, 14919556, 71650, 197972, 1087677, 204236, 11670, 19703, 82621, 32187, 16929, 21286, 2441, 6656, 2255\n        ], \n        \"k_image\": \"045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5\"\n      }\n    }\n  ], \n  \"vout\": [ {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65ab\", \n          \"view_tag\": \"bd\"\n        }\n      }\n    }, {\n      \"amount\": 0, \n      \"target\": {\n        \"tagged_key\": {\n          \"key\": \"2dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef\", \n          \"view_tag\": \"4c\"\n        }\n      }\n    }\n  ], \n  \"extra\": [ 1, 29, 234, 216, 124, 42, 64, 118, 116, 216, 13, 147, 183, 128, 67, 132, 173, 89, 54, 29, 59, 1, 98, 157, 206, 149, 126, 190, 35, 253, 189, 82, 5, 2, 9, 1, 196, 58, 221, 49, 208, 146, 193, 155\n  ], \n  \"rct_signatures\": {\n    \"type\": 6, \n    \"txnFee\": 122960000, \n    \"ecdhInfo\": [ {\n        \"trunc_amount\": \"401150d04a7559e8\"\n      }, {\n        \"trunc_amount\": \"807036489f864b37\"\n      }], \n    \"outPk\": [ \"9f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a3163195158\", \"9265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d4706\"]\n  }, \n  \"rctsig_prunable\": {\n    \"nbp\": 1, \n    \"bpp\": [ {\n        \"A\": \"f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a5772\", \n        \"A1\": \"4e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b\", \n        \"B\": \"68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc\", \n        \"r1\": \"5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100\", \n        \"s1\": \"bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00\", \n        \"d1\": \"cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c\", \n        \"L\": [ \"3aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d\", \"881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e\", \"2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd47968\", \"45d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a65464179\", \"3655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8d\", \"b2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13\", \"c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba59\"\n        ], \n        \"R\": [ \"07849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c9401\", \"8e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e30\", \"6706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d\", \"880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c7849031\", \"6e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117\", \"310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc34\", \"40cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a\"\n        ]\n      }\n    ], \n    \"CLSAGs\": [ {\n        \"s\": [ \"0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb2203809\", \"4443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a\", \"40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e\", \"7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab03\", \"7e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03\", \"d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303\", \"f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402\", \"ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0b\", \"bef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07\", \"e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170e\", \"d60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401\", \"b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d\", \"2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507\", \"920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af00501196005\", \"9614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07\", \"347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810d\"], \n        \"c1\": \"c8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205\", \n        \"D\": \"b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a\"\n      }], \n    \"pseudoOuts\": [ \"6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466\"]\n  }\n}",
				BlockHeight:     3169795,
				BlockTimestamp:  1718210909,
				Confirmations:   24,
				DoubleSpendSeen: false,
				InPool:          false,
				OutputIndices: []uint64{
					106312029,
					106312030,
				},
				PrunableAsHex: "",
				PrunableHash:  "734646af580b1e38eea27214d821670c731f92be1bea142bbb5a90b923a18292",
				PrunedAsHex:   "",
				TxHash:        "45b27c7cde61cdbb05b957241a4cc698665f7fbc55af62c84cab45f71b879a15",
				TxInfo: daemon.MoneroTxInfo{
					Version:    2,
					UnlockTime: 0,
					Vin: []daemon.Vin2{
						{
							Key: daemon.Key{
								Amount:     0,
								KeyOffsets: []int64{74944009, 14689396, 14919556, 71650, 197972, 1087677, 204236, 11670, 19703, 82621, 32187, 16929, 21286, 2441, 6656, 2255},
								KeyImage:   "045fdea0ca6f106cb9dd9da659d31af2f7f08ba79b10148a6f5d1f424d7107c5",
							},
						},
					},
					Vout: []daemon.Vout1{
						{
							Amount: 0,
							Target: daemon.Target{
								TaggedKey: daemon.TaggedKey{
									Key:     "eec278d5d419e41e67815bb71ebd7a223a89230137e98a2368c5ba2852ac65ab",
									ViewTag: "bd",
								},
							},
						},
						{
							Amount: 0,
							Target: daemon.Target{
								TaggedKey: daemon.TaggedKey{
									Key:     "2dee302749c79fef24c6129397308a35844de5710674c4ddfacc71eb00e5e1ef",
									ViewTag: "4c",
								},
							},
						},
					},
					Extra: []int32{1, 29, 234, 216, 124, 42, 64, 118, 116, 216, 13, 147, 183, 128, 67, 132, 173, 89, 54, 29, 59, 1, 98, 157, 206, 149, 126, 190, 35, 253, 189, 82, 5, 2, 9, 1, 196, 58, 221, 49, 208, 146, 193, 155},
					RctSignatures: daemon.RctSignature{
						Type:   6,
						TxnFee: 122960000,
						EcdhInfo: []daemon.EcdhInfo{
							{
								TruncAmount: "401150d04a7559e8",
							},
							{
								TruncAmount: "807036489f864b37",
							},
						},
						OutPk: []string{"9f00a45bd1925d25b9fd031f494b1aeb1e091a28d2d98d61c4391a3163195158", "9265a6b33c5476eb5820b4f0c7f2ac469041b8647436cb92ec1d53cc661d4706"},
					},
					RctsigPrunable: daemon.RctsigPrunable{
						Nbp: 1,
						Bpp: []daemon.Bpp{
							{
								A:  "f9b297c178988f2352daf8340a0fd53baf91b582c6d4aa7a8518b767577a5772",
								A1: "4e90af645534fdc1fa04e9ce4f1c4be026afb329a4f591530511131e5e0ae43b",
								B:  "68faf877b14200d7d1d9086ea29bd68d0183209edc92d394209a5191d2fe49cc",
								R1: "5b7bde842eaa416b132c7ace699f6070d7cd4fec92b8925fa413b22a3282a100",
								S1: "bd624df12f7ce700bf82a0c36844909ab76876e46195018bfbaf239697443f00",
								D1: "cf8df38c1f7ed066df4ef7670ca9db12ac431b2492759ff2ba8860313c1cca0c",
								L:  []string{"3aed2bef3e8245d565f92187de80f170b6a5466b77a27ea0eb4bcdecb2e6bb0d", "881b042fefd3fffa0578c90ce33d0b19b82e57d4a1053e2a9dc37634b2fb895e", "2f810d81002721b392c1a353772103698c5f719abd9a5b977e7f3b28ddd47968", "45d7310bc6e98f3c3cf307ad74e05046a24cd935c52d08b254bf8f6a65464179", "3655498f5eae8d665729e1cb010437bbc5f2fc35d83ea52c30c6695bd2e0bd8d", "b2d6279237322ab8889d49266ac3b55654b10c18ad626adcfe38a1a0423dda13", "c0fd48a8f8ecd5725911f5904e680d067cb40214183ede1796523ba5d083ba59"},
								R:  []string{"07849d559d2e7dceff780bf0abd17ba521ec37a4d2b67f9fea9a995c2f9c9401", "8e72c26cf10b93350fc3931a1467577a0335f7780d4c6f1b5767b34eff9a1e30", "6706c32703cbfcb633a310e010166f819c200d93799b3cb61d95366329894b6d", "880c39130b0b63b97c9cb8c564e76c0d739f798b804c30d6514dd8c3c7849031", "6e43067ec4b59f28d63eb04efd68acae0d1499c03b66c04beaf7021710106117", "310adee51daa00ba916020181e969aaa8fe223769d1c9875d88875166015bc34", "40cca9b41155b29bf643305d3b580012a431b249714991762a1725877ffe6e3a"},
							},
						},
						CLSAGs: []daemon.CLSAG{
							{
								S:  []string{"0c58b9b4eb2ec059f515eb2667947db621361e68f9dcd754135c559bb2203809", "4443e5285e74a76e999d1de61ab5dd23442995ec082e6e17bb77d2d478c62f0a", "40a241e8a3a6fd3eded5705608dcb753b020f54c5f3364d3e64628f2e5619f0e", "7a33e9d35392c211855cee063cf42e0c3242a47abe9bb94446925e72c0a4ab03", "7e5b562a34c5aed8581cd01ef7d4e9784cfc180845fe55fbd16359574644bd03", "d2c7452df728556b1e47bf3923f4084dc8c6bcbbea2006115a9138105713b303", "f3074e6874ad1dc7f7c0eddf4eb9ae3ac53a4cd15eb92783ef364a6bd86ff402", "ac2f1b59cf2eb3dc5d41f66dadb899f30a32d548adac4edbedfe1531a57d7e0b", "bef8280c7fdd5f0c571fab82c31a1eeaff2727d905cf8fd1fc0cb3acc3186d07", "e552853bc1af71068a9c4826a372bcae861a52848d65eb9c40953d53bbc6170e", "d60ebbb590744fbe887364a54d4faf8c523306c5c4f60dcf96f510b2d9567401", "b05f6e5c0e17605d5a109239df3f67731a1d996b35ea227f1e7b41630be3b00d", "2218adc2b83bf4fac6482c4d5da0b1e55470a6b491a66df9cc72e8d72ac92507", "920d03d7ca70bbe5290da49cd369bfbb6276b42bdbf02787266af00501196005", "9614e552792be8f717ca556cb1e83bd877fa5f3df055b8e4e9dbf7204fe26a07", "347357b6bd172d245121ecdf4616f0dd1332c8b3da54b5ab61fe7d71aacb810d"},
								C1: "c8c826a82de08cd2eeb458a6695742589cbb09ae5e710542cc2873dc393ac205",
								D:  "b831181566b7c10118811fde0fd71de09b4b055a84bb6f011a72df2d670e0b6a",
							},
						},
						PseudoOuts: []string{"6e6a5760dd579f55b863174753b379123d18f75aa6d3681b6ab66a9c269ff466"},
					},
				},
			},
		},
		JsonRpcFooter: defaultMoneroRpcFooter,
	}

	actual, err := test_daemon.GetTransactions(exreq.Body.TxHashes, exreq.Body.DecodeAsJson, exreq.Body.Prune, exreq.Body.Split)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expected, actual)
}

// get_version
func TestGetVersion(t *testing.T) {
	reqBody := daemon.JsonRpcGenericRequestBody[daemon.EmptyMoneroRpcParams]{JsonRpcHeader: defaultMoneroRpcHeader, Method: "get_version", Params: daemon.EmptyMoneroRpcParams{}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.EmptyMoneroRpcParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: &reqBody}
	exres := `{
	"id": "0",
	"jsonrpc": "2.0",
	"result": {
		"release": true,
		"status": "OK",
		"untrusted": false,
		"version": 196613
		}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetVersionResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetVersionResult{
			Release:       true,
			Version:       196613,
			JsonRpcFooter: defaultMoneroRpcFooter,
		},
		Error: daemon.MoneroRpcError{},
	}

	actual, err := test_daemon.GetVersion()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expected, actual)
}

// get_fee_estimate
func TestGetFeeEstimate(t *testing.T) {
	reqBody := daemon.JsonRpcGenericRequestBody[daemon.EmptyMoneroRpcParams]{JsonRpcHeader: defaultMoneroRpcHeader, Method: "get_fee_estimate", Params: daemon.EmptyMoneroRpcParams{}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.EmptyMoneroRpcParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: &reqBody}
	exres := `{
	"id": "0",
	"jsonrpc": "2.0",
	"result": {
		"credits": 0,
		"fee": 7874,
		"fees": [20000,80000,320000,4000000],
		"quantization_mask": 10000,
		"status": "OK",
		"top_hash": "",
		"untrusted": false
	}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetFeeEstimateResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetFeeEstimateResult{
			Credits:          0,
			Fee:              7874,
			Fees:             []uint64{20000, 80000, 320000, 4000000},
			QuantizationMask: 10000,
			TopHash:          "",
			JsonRpcFooter:    defaultMoneroRpcFooter,
		},
		Error: daemon.MoneroRpcError{},
	}

	actual, err := test_daemon.GetFeeEstimate()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expected, actual)
}

// get_info
func TestGetInfo(t *testing.T) {
	reqBody := daemon.JsonRpcGenericRequestBody[daemon.EmptyMoneroRpcParams]{JsonRpcHeader: defaultMoneroRpcHeader, Method: "get_info", Params: daemon.EmptyMoneroRpcParams{}}
	exreq := &daemon.MoneroRpcRequest[daemon.JsonRpcGenericRequestBody[daemon.EmptyMoneroRpcParams]]{Endpoint: daemon.DEFAULT_MONERO_RPC_ENDPOINT, Body: &reqBody}
	exres := `{
	"id": "0",
	"jsonrpc": "2.0",
	"result": {
		"adjusted_time": 1612090533,
		"alt_blocks_count": 2,
		"block_size_limit": 600000,
		"block_size_median": 300000,
		"block_weight_limit": 600000,
		"block_weight_median": 300000,
		"bootstrap_daemon_address": "",
		"busy_syncing": false,
		"credits": 0,
		"cumulative_difficulty": 86168732847545368,
		"cumulative_difficulty_top64": 0,
		"database_size": 34329849856,
		"difficulty": 225889137349,
		"difficulty_top64": 0,
		"free_space": 10795802624,
		"grey_peerlist_size": 4999,
		"height": 2286472,
		"height_without_bootstrap": 2286472,
		"incoming_connections_count": 85,
		"mainnet": true,
		"nettype": "mainnet",
		"offline": false,
		"outgoing_connections_count": 16,
		"rpc_connections_count": 1,
		"stagenet": false,
		"start_time": 1611915662,
		"status": "OK",
		"synchronized": true,
		"target": 120,
		"target_height": 2286464,
		"testnet": false,
		"top_block_hash": "b92720d8315b96e32020d04e14a0c54cc13e057d4a5beb4501be490d306fdd8f",
		"top_hash": "",
		"tx_count": 11239803,
		"tx_pool_size": 21,
		"untrusted": false,
		"update_available": false,
		"version": "0.17.1.9-release",
		"was_bootstrap_ever_used": false,
		"white_peerlist_size": 1000,
		"wide_cumulative_difficulty": "0x1322201881f9c18",
		"wide_difficulty": "0x34980ab2c5"
		}
	}`
	server := getDaemonRpcTestServer(exreq, &exres)
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &daemon.JsonRpcGenericResponse[daemon.GetInfoResult]{
		JsonRpcHeader: defaultMoneroRpcHeader,
		Result: daemon.GetInfoResult{
			Credits:                   0,
			AdjustedTime:              1612090533,
			AltBlocksCount:            2,
			BlockSizeLimit:            600000,
			BlockSizeMedian:           300000,
			BlockWeightLimit:          600000,
			BlockWeightMedian:         300000,
			BootstrapDaemonAddress:    "",
			BusySyncing:               false,
			CumulativeDifficulty:      86168732847545368,
			CumulativeDifficultyTop64: 0,
			DatabaseSize:              34329849856,
			Difficulty:                225889137349,
			DifficultyTop64:           0,
			FreeSpace:                 10795802624,
			GreyPeerlistSize:          4999,
			Height:                    2286472,
			HeightWithoutBootstrap:    2286472,
			IncomingConnectionsCount:  85,
			Mainnet:                   true,
			Nettype:                   "mainnet",
			Offline:                   false,
			OutgoingConnectionsCount:  16,
			RpcConnectionsCount:       1,
			Stagenet:                  false,
			StartTime:                 1611915662,
			Synchronized:              true,
			Target:                    120,
			TargetHeight:              2286464,
			Testnet:                   false,
			TopBlockHash:              "b92720d8315b96e32020d04e14a0c54cc13e057d4a5beb4501be490d306fdd8f",
			TopHash:                   "",
			TxCount:                   11239803,
			TxPoolSize:                21,
			UpdateAvailable:           false,
			Version:                   "0.17.1.9-release",
			WasBootstrapEverUsed:      false,
			WhitePeerlistSize:         1000,
			WideCumulativeDifficulty:  "0x1322201881f9c18",
			WideDifficulty:            "0x34980ab2c5",
			JsonRpcFooter:             defaultMoneroRpcFooter,
		},
		Error: daemon.MoneroRpcError{},
	}

	actual, err := test_daemon.GetInfo()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, expected, actual)
}
