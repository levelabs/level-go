package collection

import (
	"encoding/json"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/levelabs/level-go/common"
)

type Asset struct {
	address     ethcommon.Address `json:"address"`
	baseURI     *string           `json:"baseURI"`
	totalSupply big.Int           `json:"tokenSupply"`

	attributes *map[string]string
	traitMap   *TraitMap

	priority int64
	index    int
}

type Trait = map[string]int

type TraitMap struct {
	items map[string]int
}

type AttributeMap = map[string]*TraitMap

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

func FetchAndBuildAttributes(client Fetcher, tokenUrl string, mapping map[string]*TraitMap) error {
	res, err := client.Get(tokenUrl)
	if err != nil {
		return err
	}

	var t Token
	common.UnmarshalJSON(res, &t)

	attributes := t.Attributes

	for j := 0; j < len(attributes); j++ {
		attribute := attributes[j]
		trait := attribute.Trait
		value := attribute.Value

		if mapping[trait] == nil {
			t := new(TraitMap)
			t.items = make(map[string]int)
			mapping[trait] = t
		}

		mapping[trait].items[value]++
	}

	return nil
}
