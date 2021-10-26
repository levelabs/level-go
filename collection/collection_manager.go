package collection

import (
	"errors"
	"fmt"
	// "log"
	// "strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	// cid "github.com/ipfs/go-cid"
	// mbase "github.com/multiformats/go-multibase"
	// mh "github.com/multiformats/go-multihash"
	// shell "github.com/ipfs/go-ipfs-api"
)

var (
	ethURI  = "https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5"
	ipfsURI = "localhost:5001"

	errEthClientFailed = errors.New("failed to start eth client")
	errBaseURINotFound = errors.New("baseURI not found")
)

type CollectionManager struct {
	eth *ethclient.Client

	queue *CollectionQueue
}

func (cm *CollectionManager) RunSequence() (*Asset, error) {
	cq := cm.queue
	if cq.Len() <= 0 {
		return nil, errEmptyQueue
	}

	asset := PopCollectionQueue(cq)
	fmt.Printf("Sequencing: %.2d:%s\n", asset.priority, asset.address)

	err := cm.SetTotalSupplyForAsset(asset)
	if err != nil {
		cq.PushAndSetPriorityNow(asset)
		return nil, err
	}

	err = cm.SetBaseURIForAsset(asset)
	if err != nil {
		cq.PushAndSetPriorityNow(asset)
		return nil, err
	}

	err = cm.QueryAttributes(asset)
	if err != nil {
		cq.PushAndSetPriorityNow(asset)
		return nil, err
	}

	return asset, nil
}

func NewCollectionManager(assets map[string]int64) (*CollectionManager, error) {
	client, err := ethclient.Dial(ethURI)
	if err != nil {
		return nil, errEthClientFailed
	}

	queue := NewCollectionQueue(assets)

	cm := CollectionManager{
		eth:   client,
		queue: queue,
	}

	return &cm, nil
}

func (cm *CollectionManager) SetBaseURIForAsset(asset *Asset) error {
	collection, err := NewCollection(common.HexToAddress(asset.address), cm.eth)
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

func (cm *CollectionManager) SetTotalSupplyForAsset(asset *Asset) error {
	collection, err := NewCollection(common.HexToAddress(asset.address), cm.eth)
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

func (cm *CollectionManager) QueryAttributes(asset *Asset) error {
	if asset.baseURI == nil {
		return errBaseURINotFound
	}

	// fmt.Println("uri", *asset.baseURI)

	// sh := shell.NewShell(ipfsURI)
	// c, err := cid.Decode(strings.Trim(*asset.baseURI, "ipfs://"))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	//

	// fmt.Println(c.StringOfBase(mbase.Base58BTC))

	// fmt.Println(cid.Decode(strings.Trim(*asset.baseURI, "ipfs://")))
	// fmt.Println(mh.FromB58String(strings.Trim(*asset.baseURI, "ipfs://")))

	// fmt.Println(sh.DagGet(strings.Trim(*asset.baseURI, "ipfs://"), `{"Data": { "/": { "bytes": "string"}}}`))

	return nil
}
