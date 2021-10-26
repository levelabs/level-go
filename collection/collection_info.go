package collection

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	cid "github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
	// mh "github.com/multiformats/go-multihash"
	// shell "github.com/ipfs/go-ipfs-api"
)

var (
	ethURI  = "https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5"
	ipfsURI = "localhost:5001"

	errBaseURINotFound = errors.New("baseURI not found")
)

func (asset *Asset) SetBaseURI() error {
	client, err := ethclient.Dial(ethURI)
	if err != nil {
		log.Fatal(err)
		return err
	}

	address := common.HexToAddress(asset.address)
	collection, err := NewCollection(address, client)
	if err != nil {
		log.Fatal(err)
		return err
	}

	baseURI, err := collection.BaseURI(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	asset.baseURI = &baseURI

	return nil
}

// {"image":"ipfs://QmPbxeGcXhYQQNgsC6a36dDyYUcHgMLnGKnF8pVFmGsvqi","attributes":[{"trait_type":"Mouth","value":"Grin"},{"trait_type":"Clothes","value":"Vietnam Jacket"},{"trait_type":"Background","value":"Orange"},{"trait_type":"Eyes","value":"Blue Beams"},{"trait_type":"Fur","value":"Robot"}]}
// type Folder struct {
// }
//
// type Item struct {
// 	Hash  ItemData
// 	Name  string
// 	Tsize number
// }
//
// type ItemData struct {
// }

func (asset *Asset) QueryAttributes() error {
	if asset.baseURI == nil {
		return errBaseURINotFound
	}

	fmt.Println("uri", *asset.baseURI)

	// sh := shell.NewShell(ipfsURI)
	c, err := cid.Decode(strings.Trim(*asset.baseURI, "ipfs://"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println(c.StringOfBase(mbase.Base58BTC))

	// fmt.Println(cid.Decode(strings.Trim(*asset.baseURI, "ipfs://")))
	// fmt.Println(mh.FromB58String(strings.Trim(*asset.baseURI, "ipfs://")))

	// fmt.Println(sh.DagGet(strings.Trim(*asset.baseURI, "ipfs://"), `{"Data": { "/": { "bytes": "string"}}}`))

	return nil
}
