package collection

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
	// "log"
	"io"
	// "io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	ethURI  = "https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5"
	ipfsURI = "localhost:5001"

	errEthClientFailed  = errors.New("failed connection with eth client")
	errIPFSClientFailed = errors.New("failed connection with ipfs client")

	errBaseURINotFound   = errors.New("baseURI not found")
	errURIFormatNotFound = errors.New("An unknown URI format has been found")
)

type CollectionManager struct {
	eth  *ethclient.Client
	ipfs *shell.Shell

	queue *CollectionQueue
}

func NewCollectionManager(assets map[string]int64) (*CollectionManager, error) {
	client, err := ethclient.Dial(ethURI)
	if err != nil {
		return nil, errEthClientFailed
	}

	ipfs := shell.NewShell(ipfsURI)

	queue := NewCollectionQueue(assets)

	cm := CollectionManager{
		eth:   client,
		ipfs:  ipfs,
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

func (cm *CollectionManager) QueryTokensForAsset(asset *Asset) error {
	collection, err := NewCollection(common.HexToAddress(asset.address), cm.eth)
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

	switch baseUrl.Scheme {
	case "ipfs":
		if tokens, err := cm.ipfs.ObjectGet(baseUrl.Host); err == nil {
			for i := 0; i < len(tokens.Links); i += 1000 {
				if token, err := cm.ipfs.Cat(tokens.Links[i].Hash); err == nil {
					buf := new(strings.Builder)
					_, err := io.Copy(buf, token)
					if err != nil {
						return err
					}
					fmt.Println(buf, buf.String())
				}
			}
		}
	case "https":
		uri := strings.TrimRightFunc(uriZero, func(r rune) bool {
			return unicode.IsNumber(r)
		})
		for i := 0; i < int((asset.totalSupply).Int64()); i += 1000 {
			tokenUrl := strings.Join([]string{uri, strconv.Itoa(i)}, "")
			if res, err := http.Get(tokenUrl); err == nil {
				buf := new(strings.Builder)
				_, err := io.Copy(buf, res.Body)
				if err != nil {
					return err
				}
				fmt.Println(buf, buf.String())
			}
		}
	default:
		return errURIFormatNotFound
	}

	// ret
	return nil
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

	// err = cm.SetBaseURIForAsset(asset)
	// if err != nil {
	// 	cq.PushAndSetPriorityNow(asset)
	// 	return nil, err
	// }

	err = cm.QueryTokensForAsset(asset)
	if err != nil {
		cq.PushAndSetPriorityNow(asset)
		return nil, err
	}

	return asset, nil
}
