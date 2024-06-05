package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/monero-ecosystem/go-monero-rpc-client/daemon"
	"github.com/stretchr/testify/assert"
)

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

	test_daemon := daemon.Create(daemon.NewRpcConnection(server.URL, "", ""), 5*time.Second, &daemon.DaemonListenerHandlerAsync{})

	expected := &daemon.GetHeightResponse{
		"7e23a28cfa6df925d5b63940baf60b83c0cbb65da95f49b19e7cf0ce7dd709ce",
		2287217,
		daemon.JsonRpcFooter{"OK", false}}
	actual, err := test_daemon.GetCurrentHeight()
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, expected, actual)
}

func TestGetBlockCount(t *testing.T) {
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

	test_daemon := daemon.Create(daemon.NewRpcConnection(server.URL, "", ""), 5*time.Second, &daemon.DaemonListenerHandlerAsync{})

	expected := &daemon.MoneroRpcGenericResponse[daemon.GetBlockCountResult]{
		daemon.JsonRpcHeader{"0", "2.0"},
		daemon.GetBlockCountResult{
			24,
			daemon.JsonRpcFooter{"OK", false},
		},
	}
	actual, err := test_daemon.GetBlockCount()
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, expected, actual)
}
