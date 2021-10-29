package collection

type Item struct {
	name map[string]int
}

type Trait struct {
	Counter map[string]*Item
}

func NewTrait() *Trait {
	var trait Trait
	trait.Counter = make(map[string]*Item)
	return &trait
}

func BuildTrait(attributes *[]Attribute, trait *Trait) {
	counter := (*trait).Counter

	for j := 0; j < len(*attributes); j++ {
		trait := (*attributes)[j].Trait
		value := (*attributes)[j].Value

		if counter[trait] == nil {
			t := new(Item)
			t.name = make(map[string]int)
			counter[trait] = t
		}

		counter[trait].name[value]++
	}
}
