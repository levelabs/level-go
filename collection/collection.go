package collection

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Asset struct {
	address     common.Address `json:"address"`
	baseURI     *string        `json:"baseURI"`
	totalSupply big.Int        `json:"tokenSupply"`

	attributes *map[string]string

	priority int64
	index    int
}

type Attribute struct {
	Trait string `json:"trait_type"`
	Value string `json:"value"`
}

type Token struct {
	Image      string      `json:"image"`
	Attributes []Attribute `json:"attributes"`
}

type Trait struct {
	items map[string]int
}

func NewAsset(address string, priority int64, index int) *Asset {
	a := Asset{
		address:  common.HexToAddress(address),
		priority: priority,
		index:    index,
	}
	return &a
}

func (a *Asset) SetTotalSupply(totalSupply big.Int) {
	a.totalSupply = totalSupply
}

func (a *Asset) SetBaseURI(baseURI string) {
	a.baseURI = &baseURI
}

func (a *Asset) Address() string {
	address := a.address
	return address.String()
}

func (a *Asset) String() string {
	return fmt.Sprintf("%s - %s", a.address, (a.totalSupply).String())
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}
