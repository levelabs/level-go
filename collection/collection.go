package collection

type Collections struct {
	assets []Asset
}

type Asset struct {
	address string

	baseURI    *string
	attributes *map[string]string

	priority int64
	index    int
}

func NewAsset(address string) *Asset {
	a := Asset{address: address}
	return &a
}

func (a *Asset) Address() string {
	address := a.address
	return address
}
