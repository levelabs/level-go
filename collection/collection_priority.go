package collection

import (
	"container/heap"
	"fmt"
	"time"
)

type CollectionPriority []*Asset

func (cp CollectionPriority) Len() int { return len(cp) }

func (cp CollectionPriority) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return cp[i].priority < cp[j].priority
}

func (cp *CollectionPriority) Push(x interface{}) {
	n := len(*cp)
	asset := x.(*Asset)
	asset.index = n
	*cp = append(*cp, asset)
}

func (cp *CollectionPriority) Pop() interface{} {
	old := *cp
	n := len(old)
	asset := old[n-1]
	old[n-1] = nil   // avoid memory leak
	asset.index = -1 // for safety
	*cp = old[0 : n-1]
	return asset
}

func (cp CollectionPriority) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
	cp[i].index = i
	cp[j].index = j
}

// update modifies the priority and value of an Asset in the queue.
func (cp *CollectionPriority) update(asset *Asset, address string, priority int64) {
	asset.address = address
	asset.priority = priority
	heap.Fix(cp, asset.index)
}

func CollectionPriorityTest() {
	// Some assets and their priorities.
	assets := map[string]int64{
		"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d": time.Now().UnixNano(),
		"0x4be3223f8708ca6b30d1e8b8926cf281ec83e770": time.Now().UnixNano(),
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(CollectionPriority, len(assets))
	i := 0
	for address, priority := range assets {
		pq[i] = &Asset{
			address:  address,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new asset and then modify its priority.
	asset := &Asset{
		address:  "0x8a90cab2b38dba80c64b7734e58ee1db38b8992e",
		priority: time.Now().UnixNano(),
	}
	heap.Push(&pq, asset)

	pq.update(asset, asset.address, time.Now().UnixNano())

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		asset := heap.Pop(&pq).(*Asset)
		fmt.Printf("%.2d:%s ", asset.priority, asset.address)
	}
}
