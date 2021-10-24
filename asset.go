package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/levelabs/level-go/contracts"
)

type Attributes struct {
	traits map[string]string
}

type Asset struct {
	address    string
	baseURI    string
	attributes Attributes
}

func connect() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/79808cbe443249a8bc8bf46dea32b6f5")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")
	instance, err := contracts.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")
	_ = instance
}
