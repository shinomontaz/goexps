package utils

import (
	"testing"

	"gitlab.wildberries.ru/logisticcloud/logistic-route/lml-calculator/internal/rand"
)

// go test -bench=Heap -benchmem -benchtime=1s
// go test -bench . -benchmem -benchtime=1s

func BenchmarkHeapInsert(b *testing.B) {
	rnd := rand.New(1 << 23)
	h := NewHeap()

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		h.Insert(k, v)
	}
}

func BenchmarkHeapSelect(b *testing.B) {
	rnd := rand.New(1 << 23)
	h := NewHeap()

	var (
		k float64
		v int
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for i := 0; i < 1000; i++ {
			k = rnd.Float64()
			v = int(k * 10.0)
			h.Insert(k, v)
		}
		for i := 0; i < 1000; i++ {
			h.PopMin()
		}
	}
}

func BenchmarkHeap(b *testing.B) {
	rnd := rand.New(1 << 23)
	h := NewHeap()

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		if k > 0.5 {
			h.Insert(k, v)
		} else if !h.IsEmpty() {
			h.PopMin()
		}
	}
}

func TestHeap(t *testing.T) {
	rnd := rand.New(1 << 23)
	h := NewHeap()

	var (
		k float64
		v int
	)

	for i := 0; i < 10; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		h.Insert(k, v)
	}

	h.Clear()
	if !h.IsEmpty() {
		t.Error("Error on Clear heap: not empty!")
	}
}
