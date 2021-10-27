package collection

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type Asset struct {
	address string `json:"address"`

	baseURI     *string `json:"baseURI"`
	totalSupply big.Int `json:"tokenSupply"`

	attributes *map[string]string

	priority int64
	index    int
}

func NewAsset(address string) *Asset {
	a := Asset{address: address}
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
	return address
}

func (a *Asset) String() string {
	return fmt.Sprintf("%s - %s", a.address, (a.totalSupply).String())
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}
