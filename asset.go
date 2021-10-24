package main

import (
	// "github.com/dgraph-io/ristretto"
	"fmt"
	"net/http"
)

type Attributes struct {
	traits map[string]string
}

type Asset struct {
	address    string
	baseURI    string
	attributes Attributes
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}
