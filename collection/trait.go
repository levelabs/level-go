package collection

type Item struct {
	name map[string]int
}

type Trait struct {
	Counter map[string]*Item

	Index int
}

func NewTrait() *Trait {
	var trait Trait
	trait.Counter = make(map[string]*Item)
	return &trait
}

func NewItem() *Item {
	t := new(Item)
	t.name = make(map[string]int)
	return t
}

func BuildTrait(attributes *[]Attribute, trait *Trait) {
	counter := (*trait).Counter

	for j := 0; j < len(*attributes); j++ {
		trait := (*attributes)[j].Trait
		value := (*attributes)[j].Value

		if counter[trait] == nil {
			item := NewItem()
			counter[trait] = item
		}

		counter[trait].name[value]++
	}

	(*trait).Index++
}
