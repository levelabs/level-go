package collection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	ethURI  = "https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5"
	ipfsURI = "localhost:5001"

	errClientEthereumFailed = errors.New("failed connection with eth client")
	errClientIPFSFailed     = errors.New("failed connection with ipfs client")

	errBaseURINotFound   = errors.New("baseURI not found")
	errURIFormatNotFound = errors.New("An unknown URI format has been found")

	errEmptyWaitlist = errors.New("Waitlist is empty")
)

type Client struct {
	Ethereum *ethclient.Client
	IPFS     *shell.Shell
}

type Manager struct {
	Connection *Client
	Waitlist   *PriorityQueue
}

func NewManager(assets map[string]int64) (*Manager, error) {
	ethereum, err := ethclient.Dial(ethURI)
	if err != nil {
		return nil, errClientEthereumFailed
	}

	ipfs := shell.NewShell(ipfsURI)
	waitlist := NewPriorityQueue(assets)

	connection := &Client{
		Ethereum: ethereum,
		IPFS:     ipfs,
	}

	manager := Manager{
		Connection: connection,
		Waitlist:   waitlist,
	}

	return &manager, nil
}

func (manager *Manager) SetBaseURIForAsset(asset *Asset) error {
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum)
	if err != nil {
		return err
	}

	baseURI, err := collection.BaseURI(&bind.CallOpts{})
	if err != nil {
		return err
	}

	asset.SetBaseURI(baseURI)

	return nil
}

func (manager *Manager) SetTotalSupplyForAsset(asset *Asset) error {
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum)
	if err != nil {
		return err
	}

	totalSupply, err := collection.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	asset.SetTotalSupply(*totalSupply)

	return nil
}

func (manager *Manager) QueryTokensForAsset(asset *Asset) error {
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum)
	if err != nil {
		return err
	}

	// todo: what if there isn't a token zero
	tokenZero, err := collection.TokenByIndex(&bind.CallOpts{}, big.NewInt(0))
	if err != nil {
		return err
	}

	uriZero, err := collection.TokenURI(&bind.CallOpts{}, tokenZero)
	if err != nil {
		return err
	}

	baseUrl, err := url.Parse(uriZero)
	if err != nil {
		return err
	}

	attributes := make(map[string]*Trait)
	switch baseUrl.Scheme {
	case "ipfs":
		if ipfsUris, err := manager.Connection.IPFS.ObjectGet(baseUrl.Host); err == nil {
			for i := 0; i < len(ipfsUris.Links); i += 1000 {
				if token, err := manager.Connection.IPFS.Cat(ipfsUris.Links[i].Hash); err == nil {
					var t Token
					buf, _ := ioutil.ReadAll(token)
					json.Unmarshal(buf, &t)

					for j := 0; j < len(t.Attributes); j++ {
						attribute := t.Attributes[j]
						trait := attribute.Trait
						value := attribute.Value

						if attributes[trait] == nil {
							var t *Trait
							t = new(Trait)
							t.items = make(map[string]int)
							attributes[trait] = t
						}

						attributes[trait].items[value]++
					}
				}
			}
		}
	case "https":
		for i := 0; i < int((asset.totalSupply).Int64()); i += 1 {
			tokenUrl := strings.Join([]string{strings.TrimRightFunc(uriZero, func(r rune) bool {
				return unicode.IsNumber(r)
			}), strconv.Itoa(i)}, "")
			if res, err := http.Get(tokenUrl); err == nil {
				var t Token
				buf, _ := ioutil.ReadAll(res.Body)
				json.Unmarshal(buf, &t)

				for j := 0; j < len(t.Attributes); j++ {
					attribute := t.Attributes[j]
					trait := attribute.Trait
					value := attribute.Value

					if attributes[trait] == nil {
						var t *Trait
						t = new(Trait)
						t.items = make(map[string]int)
						attributes[trait] = t
					}

					attributes[trait].items[value]++
				}
			}
		}
	default:
		return errURIFormatNotFound
	}

	fmt.Println(attributes["Background"])
	fmt.Println(attributes["Fur"])
	fmt.Println(attributes["Clothes"])

	// ret
	return nil
}

func (manager *Manager) RunSequence() (*Asset, error) {
	if manager.Waitlist.Len() <= 0 {
		return nil, errEmptyWaitlist
	}

	asset, err := manager.WaitlistRemove()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Sequencing: %.2d:%s\n", asset.priority, asset.address)

	err = manager.SetTotalSupplyForAsset(asset)
	if err != nil {
		manager.WaitlistAppend(asset)
		return nil, err
	}

	err = manager.QueryTokensForAsset(asset)
	if err != nil {
		manager.WaitlistAppend(asset)
		return nil, err
	}

	return asset, nil
}

func (manager *Manager) WaitlistAppend(asset *Asset) {
	manager.Waitlist.PriorityQueuePush(asset)
}

func (manager *Manager) WaitlistRemove() (*Asset, error) {
	asset, err := manager.Waitlist.PriorityQueueRemove()
	if err != nil {
		return nil, err
	}
	return asset, nil
}
