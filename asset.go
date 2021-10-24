package main

import (
// "github.com/dgraph-io/ristretto"
// "fmt"
)

type Attributes struct {
	traits map[string]string
}

type Asset struct {
	address    string
	baseURI    string
	attributes Attributes
}
