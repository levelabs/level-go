package collection

import (
	"container/heap"
	"errors"
	"fmt"
	"time"
)

var (
	errEmptyQueue = errors.New("collection queue is empty")
)

type CollectionQueue []*Asset

func (cq CollectionQueue) Len() int { return len(cq) }

func (cq CollectionQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return cq[i].priority < cq[j].priority
}

func (cq *CollectionQueue) Push(x interface{}) {
	n := len(*cq)
	asset := x.(*Asset)
	asset.index = n
	*cq = append(*cq, asset)
}

func (cq *CollectionQueue) Pop() interface{} {
	old := *cq
	n := len(old)
	asset := old[n-1]
	old[n-1] = nil   // avoid memory leak
	asset.index = -1 // for safety
	*cq = old[0 : n-1]
	return asset
}

func (cq CollectionQueue) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
	cq[i].index = i
	cq[j].index = j
}

// update modifies the priority and value of an Asset in the queue.
func (cq *CollectionQueue) update(asset *Asset, priority int64) {
	asset.priority = priority
	heap.Fix(cq, asset.index)
}

func NewCollectionQueue(assets map[string]int64) *CollectionQueue {
	cq := make(CollectionQueue, len(assets))
	i := 0
	for address, priority := range assets {
		cq[i] = &Asset{
			address:  address,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&cq)
	return &cq
}

func PopCollectionQueue(cq *CollectionQueue) *Asset {
	asset := heap.Pop(cq).(*Asset)
	return asset
}

func CollectionQueueTest() {
	// Some assets and their priorities.
	assets := map[string]int64{
		"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d": time.Now().UnixNano(),
		"0x4be3223f8708ca6b30d1e8b8926cf281ec83e770": time.Now().UnixNano(),
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	cq := NewCollectionQueue(assets)

	// Insert a new asset and then modify its priority.
	asset := &Asset{
		address:  "0x8a90cab2b38dba80c64b7734e58ee1db38b8992e",
		priority: time.Now().UnixNano(),
	}
	heap.Push(cq, asset)

	cq.update(asset, time.Now().UnixNano())

	// Take the items out; they arrive in decreasing priority order.
	for cq.Len() > 0 {
		asset := heap.Pop(cq).(*Asset)
		fmt.Printf("%.2d:%s ", asset.priority, asset.address)
	}
}
