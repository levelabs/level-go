package collection

import (
	"errors"
	"io"
	net "net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	shell "github.com/ipfs/go-ipfs-api"
)

var (
	errClientEthereumFailed = errors.New("failed connection with eth client")
	errClientIPFSFailed     = errors.New("failed connection with ipfs client")

	errClientIPFSGet = errors.New("There was an issue when querying the IPFS client")
)

type ClientFetcher interface {
	Get(uri string) (io.ReadCloser, error)
}

type Client struct {
	Ethereum Ethereum
	IPFS     IPFS
	Http     Http
}

type IPFS struct {
	Client *shell.Shell
}

type Ethereum struct {
	Client *ethclient.Client
}

type Http struct {
	Client net.Client
}

type ClientConfig struct {
	EthUri  string
	IPFSUri string
}

func BuildClient(config ClientConfig) (*Client, error) {
	ipfsUri := config.IPFSUri
	ethUri := config.EthUri

	ipfs := IPFS{
		Client: shell.NewShell(ipfsUri),
	}

	eth, err := ethclient.Dial(ethUri)
	if err != nil {
		return nil, errClientEthereumFailed
	}

	ethereum := Ethereum{
		Client: eth,
	}

	http := Http{
		Client: net.Client{},
	}

	client := &Client{
		Ethereum: ethereum,
		IPFS:     ipfs,
		Http:     http,
	}

	return client, nil
}

func (ipfs IPFS) Get(uri string) (io.ReadCloser, error) {
	res, err := ipfs.Client.Cat(uri)
	if err != nil {
		return nil, errClientIPFSGet
	}
	return res, nil
}

func (http Http) Get(uri string) (io.ReadCloser, error) {
	res, err := http.Client.Get(uri)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
