package wallet

import (
	"net/http"
)

// Config holds the configuration of a monero rpc client.
type Config struct {
	Address       string
	CustomHeaders map[string]string
	Transport     http.RoundTripper
}
