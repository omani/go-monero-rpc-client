package test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/monero-ecosystem/go-monero-rpc-client/daemon"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func getUrlFromEnv(name string) (*url.URL, error) {
	urlStr := strings.TrimSpace(os.Getenv(name))
	if urlStr == "" {
		return nil, errors.New(name + " env can't be empty")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func initSimpleDaemonRpcClient(t *testing.T) daemon.IDaemonRpcClient {
	return initDaemonRpcClientWithCreds(t, "", "")
}

func initDaemonRpcClientWithCreds(t *testing.T, username, password string) daemon.IDaemonRpcClient {
	u, err := getUrlFromEnv("MONERO_DAEMON_RPC_ADDRESS")
	if err != nil {
		t.Fatal(err.Error())
	}

	return daemon.NewDaemonRpcClient(daemon.NewRpcConnection(u, username, password))
}

func initDaemonRpcClientWithCredsAndCustomUrl(t *testing.T, urlStr, username, password string) daemon.IDaemonRpcClient {
	if urlStr == "" {
		t.Fatal("Url can't be empty")
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	return daemon.NewDaemonRpcClient(daemon.NewRpcConnection(u, username, password))
}

func TestGetTransactionPool(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetTransactionPool()
	if err != nil {
		t.Error(err.Error())
	}
}
func TestGetBlockByHash(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockByHash(false, "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetBlockByheight(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockByHeight(false, 2751506)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetBlockCount(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockCount()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetBlockHeaderByHash(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockHeaderByHash(false, "43bd1f2b6556dcafa413d8372974af59e4e8f37dbf74dc6b2a9b7212d0577428")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetBlockHeaderByHe(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockHeaderByHeight(false, 2751506)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetBlockHeadersRange(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockHeadersRange(false, 2751506, 2751507)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetBlockTemplate(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetBlockTemplate("48ukkZtBSBRL8iva7k3p2sBVMLWTfNwsTbW1aVh5M84g21muDCssvCHTpoZCaSc6rq8M9QLZ3sQMrMn1bq2RD2anGnyHhtq", 123)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetCurrentHeight(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetCurrentHeight()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetLastBlockHeader(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.GetLastBlockHeader(false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestOnGetBlockHash(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	_, err := c.OnGetBlockHash(2751506)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTransactions(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)

	cases := []struct {
		p1 []string
		p2 bool
		p3 bool
		p4 bool
	}{
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, true, false, false},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, true, true, false},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, false, false, true},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, true, true, true},
		{[]string{"260823575cf078f999eb9d9355892790b0aa11bf33c580203c57dc514ea0e70b"}, false, false, false},
	}

	wait := &sync.WaitGroup{}
	wait.Add(len(cases))
	errorChan := make(chan error, len(cases))

	for _, v := range cases {
		go func(params struct {
			p1 []string
			p2 bool
			p3 bool
			p4 bool
		}) {
			_, err := c.GetTransactions(params.p1, params.p2, params.p3, params.p4)
			if err != nil {
				errorChan <- err
			}
		}(v)
	}
	wait.Done()
	close(errorChan)

	for v := range errorChan {
		t.Error(v.Error())
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)
	_, err := c.GetVersion()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetFeeEstimate(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)
	_, err := c.GetFeeEstimate()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetInfo(t *testing.T) {
	t.Parallel()

	c := initSimpleDaemonRpcClient(t)
	_, err := c.GetInfo()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDigestAuth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	containerReq := testcontainers.ContainerRequest{
		Image:        "sethsimmons/simple-monerod:v0.18.3.3",
		ExposedPorts: []string{"18081/tcp"},
		Cmd:          []string{"--rpc-restricted-bind-ip=0.0.0.0", "--rpc-bind-ip=0.0.0.0", "--confirm-external-bind", "--rpc-login=user:pass", "--offline"},
		WaitingFor:   wait.ForLog(`Use "help <command>" to see a command's documentation.`),
	}

	monerodC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	defer monerodC.Terminate(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	u, err := monerodC.PortEndpoint(ctx, "18081/tcp", "http")
	if err != nil {
		t.Fatal(err.Error())
	}

	client := initDaemonRpcClientWithCredsAndCustomUrl(t, u, "user", "pass")

	_, err = client.GetBlockCount()
	if err != nil {
		t.Error(err.Error())
	}

}
