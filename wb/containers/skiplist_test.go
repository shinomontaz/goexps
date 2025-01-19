package utils

import (
	"testing"

	"gitlab.wildberries.ru/logisticcloud/logistic-route/lml-calculator/internal/rand"
)

// go test -bench=Skiplist -benchmem -benchtime=1s
// go test -bench . -benchmem -benchtime=1s

func BenchmarkSkiplistInsert(b *testing.B) {
	rnd := rand.New(1 << 23)
	sl := NewSkipList(WithRnd(rnd))

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		sl.Insert(k, v)
	}
}

func BenchmarkSkiplist(b *testing.B) {
	rnd := rand.New(1 << 23)
	sl := NewSkipList(WithRnd(rnd))

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		if k > 0.5 {
			sl.Insert(k, v)
		} else {
			sl.PopMin()
		}
	}
}

func BenchmarkSkiplistSelect(b *testing.B) {
	rnd := rand.New(1 << 23)
	sl := NewSkipList(WithRnd(rnd))

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//		b.StopTimer()
		for i := 0; i < 1000; i++ {
			k = rnd.Float64()
			v = int(k * 10.0)
			sl.Insert(k, v)
		}
		//		b.StartTimer()
		for i := 0; i < 1000; i++ {
			sl.PopMin()
		}
	}
}

func TestSkiplist(t *testing.T) {
	rnd := rand.New(1 << 23)
	sl := NewSkipList(WithRnd(rnd))

	var (
		k float64
		v int
	)

	for i := 0; i < 10; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		sl.Insert(k, v)
	}
	prev := 0.0
	for i := 0; i < 10; i++ {
		k, _, err := sl.PopMin()
		if err != nil {
			t.Errorf("Error on PopMin: %v", err)
		}

		if prev > k {
			t.Errorf("Error on PopMin - wrong order: got %f, prev: %v", k, prev)
		}

		prev = k
	}
}
