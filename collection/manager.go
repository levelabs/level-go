package collection

import (
	"errors"
	"log"

	"github.com/levelabs/level-go/common"
)

const (
	ethUri  = "https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5"
	ipfsUri = "localhost:5001"
)

var (
	errEmptyWaitlist     = errors.New("Waitlist is empty")
	errAttributesUpdated = errors.New("Attributes have been updated")
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

func NewManager(
	assets map[string]int64,
) (*Manager, error) {
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

func (manager *Manager) RunSequence() (*Asset, error) {
	if manager.Waitlist.Len() <= 0 {
		return nil, errEmptyWaitlist
	}

	asset, err := manager.WaitlistRemove()
	if err != nil {
		return nil, err
	}
	log.Printf("[SEQUENCE]: %s:%.2d\n", asset.address, asset.priority)

	trait, err := manager.UpdateAttributes(asset)
	if err != nil {
		manager.WaitlistAppend(asset)
		return nil, err
	}

	asset.trait = trait

	return asset, nil
}

func (manager *Manager) UpdateAttributes(asset *Asset) (*Trait, error) {
	err := asset.SetBaseUri(&manager.Connection.Ethereum)
	if err != nil {
		// Handle Errors Here!
		return nil, err
	}

	trait := NewTrait()

	// todo: add arweave getter
	switch asset.uri.Scheme {
	case UriIPFS:
		if err := manager.RunIPFSTraitGetter(trait, asset); err != nil {
			return nil, err
		}
	case UriHttp:
		if err := manager.RunHttpTraitGetter(trait, asset); err != nil {
			return nil, err
		}
	default:
		return nil, errURIFormatNotFound
	}

	// ret
	return trait, nil
}

func (manager *Manager) RunIPFSTraitGetter(trait *Trait, asset *Asset) error {
	ipfsUris, err := manager.Connection.IPFS.Client.ObjectGet(asset.uri.Host)
	if err != nil {
		return err
	}

	for i := 0; i < len(ipfsUris.Links); i += 5000 {
		var token Token
		err := GetTokenData(manager.Connection.IPFS, ipfsUris.Links[i].Hash, &token)
		if err != nil {
			return err
		}
		BuildTrait(&token.Attributes, trait)
	}

	return nil
}

// todo: fix http getter
func (manager *Manager) RunHttpTraitGetter(trait *Trait, asset *Asset) error {
	for i := 0; i < int((asset.totalSupply).Int64()); i += 5000 {
		var token Token
		err := GetTokenData(manager.Connection.Http, common.BuildUrl(asset.uri.Host, i), &token)
		if err != nil {
			return err
		}
		BuildTrait(&token.Attributes, trait)
	}
	return nil
}

func GetTokenData(fetcher ClientFetcher, tokenUrl string, token *Token) error {
	res, err := fetcher.Get(tokenUrl)
	if err != nil {
		return err
	}
	if err := common.UnmarshalJSON(res, &token); err != nil {
		return err
	}
	return nil
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
