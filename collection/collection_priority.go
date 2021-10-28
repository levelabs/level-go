package collection

import (
	"container/heap"
	"time"
)

type PriorityQueue []*Asset

func (cq PriorityQueue) Len() int { return len(cq) }

func (cq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return cq[i].priority < cq[j].priority
}

func (cq *PriorityQueue) Push(x interface{}) {
	n := len(*cq)
	asset := x.(*Asset)
	asset.index = n
	*cq = append(*cq, asset)
}

func (cq *PriorityQueue) Pop() interface{} {
	old := *cq
	n := len(old)
	asset := old[n-1]
	old[n-1] = nil   // avoid memory leak
	asset.index = -1 // for safety
	*cq = old[0 : n-1]
	return asset
}

func (cq PriorityQueue) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
	cq[i].index = i
	cq[j].index = j
}

// update modifies the priority and value of an Asset in the PriorityQueue.
func (cq *PriorityQueue) update(asset *Asset, priority int64) {
	asset.priority = priority
	heap.Fix(cq, asset.index)
}

func NewPriorityQueue(assets map[string]int64) *PriorityQueue {
	cq := make(PriorityQueue, len(assets))
	i := 0
	for address, priority := range assets {
		cq[i] = NewAsset(address, priority, i)
		i++
	}
	heap.Init(&cq)
	return &cq
}

func (pq *PriorityQueue) PriorityQueuePush(asset *Asset) {
	heap.Push(pq, asset)
	pq.update(asset, time.Now().UnixNano())
}

func (pq *PriorityQueue) PriorityQueueRemove() (*Asset, error) {
	asset := heap.Pop(pq).(*Asset)
	return asset, nil
}
