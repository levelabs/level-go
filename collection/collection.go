package collection

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/levelabs/level-go/common"
)

var (
	errFailedConversionToBytes      = errors.New("Failed to convert asset to byte[]")
	errURIFormatNotFound            = errors.New("An unknown URI format has been found")
	errCreatingCollectionEthBinding = errors.New("There was an issue creating the eth binding ")
	errBaseUriNotExist              = errors.New("Couldn't find the base uri")
)

const (
	UriIPFS    = 1
	UriHttp    = 2
	UriArweave = 3
)

type Uri struct {
	Scheme int
	Host   string
}

// baseURI     *string           `json:"baseURI"`
type Asset struct {
	address     ethcommon.Address `json:"address"`
	uri         *Uri              `json:"uri"`
	totalSupply big.Int           `json:"tokenSupply"`

	trait *Trait

	priority int64
	index    int
}

func NewAsset(address string, priority int64, index int) *Asset {
	a := Asset{
		address:  ethcommon.HexToAddress(address),
		priority: priority,
		index:    index,
	}
	return &a
}

func (a *Asset) SetTotalSupply(totalSupply big.Int) {
	a.totalSupply = totalSupply
}

func (a *Asset) Address() string {
	return a.address.String()
}

func (a *Asset) AddressBytes() []byte {
	return a.address.Bytes()
}

func (a *Asset) String() string {
	return fmt.Sprintf("%s - %s", a.address, (a.totalSupply).String())
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(a.String())
	if err != nil {
		return nil, errFailedConversionToBytes
	}
	return bytes, nil
}

func (a *Asset) SetBaseUri(ethereum *Ethereum) error {
	collection, err := NewCollection(a.address, ethereum.Client)
	if err != nil {
		return errCreatingCollectionEthBinding
	}

	uri, err := collection.BaseURI(&bind.CallOpts{})
	if err != nil {
		return errBaseUriNotExist
	}

	base, err := url.Parse(uri)
	if err != nil {
		return err
	}

	if a.uri == nil {
		var u Uri
		switch base.Scheme {
		case "ipfs":
			u.Scheme = UriIPFS
			u.Host = base.Host
		case "https":
			u.Scheme = UriHttp
			u.Host = common.TrimRightNumber(uri)
		default:
			return errURIFormatNotFound
		}
		a.uri = &u
	}
	return nil
}

// func (a *Asset) RandomTokenBaseUri(collection *Collection) (*string, error) {
// 	tokenIndex := big.NewInt(0) // using index 0
//
// 	// todo: what if there isn't a token zero
// 	tokenId, err := collection.TokenByIndex(&bind.CallOpts{}, tokenIndex)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	uri, err := collection.TokenURI(&bind.CallOpts{}, tokenId)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return uri, nil
// }
