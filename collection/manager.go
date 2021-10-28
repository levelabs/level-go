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

func (manager *Manager) UpdateAttributes(asset *Asset) error {
	collection, err := NewCollection(asset.address, manager.Connection.Ethereum.Client)
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

	attributes := make(TraitMap)

	switch baseUrl.Scheme {
	case "ipfs":
		ipfsUris, err := manager.Connection.IPFS.Client.ObjectGet(baseUrl.Host)
		if err != nil {
			return err
		}

		for i := 0; i < len(ipfsUris.Links); i += 1 {
			err := FetchAndBuildAttributes(manager.Connection.IPFS, ipfsUris.Links[i].Hash, attributes)
			if err != nil {
				return err
			}
		}
	case "https":
		for i := 0; i < int((asset.totalSupply).Int64()); i += 1000 {
			err := FetchAndBuildAttributes(manager.Connection.Http, common.BuildUrl(uriZero, i), attributes)
			if err != nil {
				return err
			}
		}
	default:
		return errURIFormatNotFound
	}

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

	err = manager.UpdateAttributes(asset)
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
