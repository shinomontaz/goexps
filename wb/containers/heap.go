package utils

import (
	"container/heap"
)

type Heap []item

func NewHeap() *Heap {
	h := &Heap{}
	heap.Init(h)
	return h
}

func (h *Heap) PopMin() int {
	return heap.Pop(h).(int)
}

func (h *Heap) Insert(k float64, v int) {
	heap.Push(h, item{i: v, w: k})
}

func (h Heap) IsEmpty() bool {
	return len(h) == 0
}

func (h *Heap) Clear() {
	*h = (*h)[:0]
}

func (h Heap) Less(i, j int) bool {
	return h[i].w < h[j].w
}

func (h Heap) Len() int { return len(h) }

func (h Heap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(i interface{}) {
	*h = append(*h, i.(item))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	i := old[n-1]
	*h = old[0 : n-1]
	return i.i
}
