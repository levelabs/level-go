package collection

import (
	"errors"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/levelabs/level-go/common"
)

const (
	ethUri  = "https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5"
	ipfsUri = "localhost:5001"
)

var (
	errBaseURINotFound   = errors.New("baseURI not found")
	errURIFormatNotFound = errors.New("An unknown URI format has been found")
	errEmptyWaitlist     = errors.New("Waitlist is empty")
)

type Manager struct {
	Connection *Client
	Waitlist   *PriorityQueue
}

type Attribute struct {
	Trait string `json:"trait_type"`
	Value string `json:"value"`
}

type Token struct {
	Image      string      `json:"image"`
	Attributes []Attribute `json:"attributes"`
}

func NewManager(assets map[string]int64) (*Manager, error) {
	clientConfig := ClientConfig{
		EthUri:  ethUri,
		IPFSUri: ipfsUri,
	}

	client, err := BuildClient(clientConfig)
	if err != nil {
		// todo: should fail
		return nil, err
	}

	waitlist := NewPriorityQueue(assets)

	manager := Manager{
		Connection: client,
		Waitlist:   waitlist,
	}

	return &manager, nil
}

func (manager *Manager) SetBaseURIForAsset(asset *Asset) error {
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum.Client)
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
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum.Client)
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

func (manager *Manager) UpdateAttributes(asset *Asset) (*AttributeMap, error) {
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum.Client)
	if err != nil {
		return nil, err
	}

	if err := manager.SetTotalSupplyForAsset(asset); err != nil {
		return nil, err
	}

	// todo: what if there isn't a token zero
	tokenZero, err := collection.TokenByIndex(&bind.CallOpts{}, big.NewInt(0))
	if err != nil {
		return nil, err
	}

	uriZero, err := collection.TokenURI(&bind.CallOpts{}, tokenZero)
	if err != nil {
		return nil, err
	}

	baseUrl, err := url.Parse(uriZero)
	if err != nil {
		return nil, err
	}

	attributeMap := make(AttributeMap)

	switch baseUrl.Scheme {
	case "ipfs":
		ipfsUris, err := manager.Connection.IPFS.Client.ObjectGet(baseUrl.Host)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(ipfsUris.Links); i += 1000 {
			if err := FetchAndBuildAttributes(manager.Connection.IPFS, ipfsUris.Links[i].Hash, attributeMap); err != nil {
				return nil, err
			}
		}
	case "https":
		for i := 0; i < int((asset.totalSupply).Int64()); i += 1000 {
			if err := FetchAndBuildAttributes(manager.Connection.Http, common.BuildUrl(uriZero, i), attributeMap); err != nil {
				return nil, err
			}
		}
	default:
		return nil, errURIFormatNotFound
	}

	// ret
	return &attributeMap, nil
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

	if _, err := manager.UpdateAttributes(asset); err != nil {
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
