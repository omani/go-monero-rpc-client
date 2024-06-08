package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/monero-ecosystem/go-monero-rpc-client/daemon"
	"github.com/stretchr/testify/assert"
)

var defaultMoneroRpcHeader = daemon.JsonRpcHeader{"0", "2.0"}

func createTestDaemonRpcClient(u string) (*daemon.DaemonRpcClient, error) {
	u1, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return daemon.CreateDaemonRpcClient(daemon.NewRpcConnection(*u1, "", ""), 5*time.Second, &daemon.DaemonListenerHandlerAsync{}), nil
}

func daemonRpcTestServerCheck[P daemon.MoneroRpcRequestParams](r *http.Request, expected *daemon.MoneroRpcRequest[P]) bool {
	if r == nil {
		return false
	}

	if r.URL.Path != "/json_rpc" || r.Method != http.MethodPost || r.Header.Get("Accept") != "application/json" {
		return false
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	var actual daemon.MoneroRpcRequest[P]
	if err := json.Unmarshal(data, &actual); err != nil {
		return false
	}

	return reflect.DeepEqual(&actual, expected)
}
func getDaemonRpcTestServer[P daemon.MoneroRpcRequestParams](req *daemon.MoneroRpcRequest[P], res *string) *httptest.Server {
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/get_height" && r.Method == http.MethodGet && r.Header.Get("Accept") == "application/json" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"hash": "7e23a28cfa6df925d5b63940baf60b83c0cbb65da95f49b19e7cf0ce7dd709ce",
				"height": 2287217,
				"status": "OK",
				"untrusted": false
			  }`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	test_daemon, err := createTestDaemonRpcClient(server.URL)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	expected := &daemon.GetHeightResponse{
		"7e23a28cfa6df925d5b63940baf60b83c0cbb65da95f49b19e7cf0ce7dd709ce",
		2287217,
		daemon.MoneroRpcError{},
		daemon.JsonRpcFooter{"OK", false}}
	actual, err := test_daemon.GetCurrentHeight()
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockCount(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.GetBlockCountParams]{defaultMoneroRpcHeader, "get_block_count", daemon.GetBlockCountParams{}}
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
		return
	}

	expected := &daemon.MoneroRpcGenericResponse[daemon.GetBlockCountResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetBlockCountResult{
			993163,
			daemon.JsonRpcFooter{"OK", false},
		},
		daemon.MoneroRpcError{},
	}
	actual, err := test_daemon.GetBlockCount()
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestOnGetBlockHash(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.OnGetBlockHashParams]{defaultMoneroRpcHeader, "on_get_block_hash", [1]uint64{912345}}
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
		return
	}

	expected := &daemon.MoneroRpcGenericResponse[daemon.OnGetBlockHashResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		"e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6",
		daemon.MoneroRpcError{}}

	actual, err := test_daemon.OnGetBlockHash(exreq.Params[0])
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockTemplate(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.GetBlockTemplateParams]{
		defaultMoneroRpcHeader,
		"get_block_template",
		daemon.GetBlockTemplateParams{
			"44GBHzv6ZyQdJkjqZje6KLZ3xSyN1hBSFAnLP6EAqJtCRVzMzZmeXTC2AHKDS9aEDTRKmo6a6o9r9j86pYfhCWDkKjbtcns", 60}}

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
		return
	}

	expected := &daemon.MoneroRpcGenericResponse[daemon.GetBlockTemplateResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetBlockTemplateResult{
			"0e0ed286da8006ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd00000000d130d22cf308b308498bbc16e2e955e7dbd691e6a8fab805f98ad82e6faa8bcc06",
			"0e0ed286da8006ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd0000000002abc78b0101ffefc68b0101fcfcf0d4b422025014bb4a1eade6622fd781cb1063381cad396efa69719b41aa28b4fce8c7ad4b5f019ce1dc670456b24a5e03c2d9058a2df10fec779e2579753b1847b74ee644f16b023c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000051399a1bc46a846474f5b33db24eae173a26393b976054ee14f9feefe99925233802867097564c9db7a36af5bb5ed33ab46e63092bd8d32cef121608c3258edd55562812e21cc7e3ac73045745a72f7d74581d9a0849d6f30e8b2923171253e864f4e9ddea3acb5bc755f1c4a878130a70c26297540bc0b7a57affb6b35c1f03d8dbd54ece8457531f8cba15bb74516779c01193e212050423020e45aa2c15dcb",
			226807339040,
			0,
			1182367759996,
			2286447,
			"",
			"ecdc1aab3033cf1716c52f13f9d8ae0051615a2453643de94643b550d543becd",
			130,
			"d432f499205150873b2572b5f033c9c6e4b7c6f3394bd2dd93822cd7085e7307",
			2285568,
			"0x34cec55820",
			daemon.JsonRpcFooter{"OK", false}},
		daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockTemplate(exreq.Params.WalletAddress, exreq.Params.ReserveSize)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetLastBlockHeader(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.GetLastBlockHeaderParams]{
		defaultMoneroRpcHeader,
		"get_last_block_header",
		daemon.GetLastBlockHeaderParams{}}

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
		return
	}

	expected := &daemon.MoneroRpcGenericResponse[daemon.GetLastBlockHeaderResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetLastBlockHeaderResult{
			daemon.BlockHeader{
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
			0,
			"",
			daemon.JsonRpcFooter{"OK", false}},
		daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetLastBlockHeader(false)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetLastBlockHeaderByHash(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.GetBlockHeaderByHashParams]{
		defaultMoneroRpcHeader,
		"get_block_header_by_hash",
		daemon.GetBlockHeaderByHashParams{daemon.GetLastBlockHeaderParams{false}, "e22cf75f39ae720e8b71b3d120a5ac03f0db50bba6379e2850975b4859190bc6"}}

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
		return
	}

	expected := &daemon.MoneroRpcGenericResponse[daemon.GetBlockHeaderByHashResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetBlockHeaderByHashResult{
			daemon.BlockHeader{
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
			0,
			"",
			daemon.JsonRpcFooter{"OK", false}},
		daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockHeaderByHash(exreq.Params.FillPowHash, exreq.Params.Hash)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetLastBlockHeaderByHeight(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.GetBlockHeaderByHeightParams]{
		defaultMoneroRpcHeader,
		"get_block_header_by_height",
		daemon.GetBlockHeaderByHeightParams{daemon.GetLastBlockHeaderParams{false}, 912345}}

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
		return
	}

	expected := &daemon.MoneroRpcGenericResponse[daemon.GetBlockHeaderByHeightResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetBlockHeaderByHeightResult{
			daemon.BlockHeader{
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
			0,
			"",
			daemon.JsonRpcFooter{"OK", false}},
		daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockHeaderByHeight(exreq.Params.FillPowHash, exreq.Params.Height)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockHeadersRange(t *testing.T) {
	exreq := &daemon.MoneroRpcRequest[daemon.GetBlockHeadersRangeParams]{
		defaultMoneroRpcHeader,
		"get_block_headers_range",
		daemon.GetBlockHeadersRangeParams{daemon.GetLastBlockHeaderParams{false}, 1545999, 1546000}}

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
		return
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
	expected := &daemon.MoneroRpcGenericResponse[daemon.GetBlockHeadersRangeResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetBlockHeadersRangeResult{
			headers,
			0,
			"",
			daemon.JsonRpcFooter{"OK", false}},
		daemon.MoneroRpcError{}}

	actual, err := test_daemon.GetBlockHeadersRange(exreq.Params.FillPowHash, exreq.Params.StartHeight, exreq.Params.EndHeight)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, expected, actual)
}
