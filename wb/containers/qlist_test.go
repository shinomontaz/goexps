package utils

import (
	"testing"

	"gitlab.wildberries.ru/logisticcloud/logistic-route/lml-calculator/internal/rand"
)

// go test -bench=Qlist -benchmem -benchtime=1s
// go test -bench . -benchmem -benchtime=1s

func BenchmarkQlistInsert(b *testing.B) {
	rnd := rand.New(1 << 23)
	ql := NewQlist()

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		ql.Insert(k, v)
	}
}

func BenchmarkQlistSelect(b *testing.B) {
	rnd := rand.New(1 << 23)
	ql := NewQlist()

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
			ql.Insert(k, v)
		}
		//		b.StartTimer()
		for i := 0; i < 1000; i++ {
			ql.PopMin()
		}
	}
}

func BenchmarkQlist(b *testing.B) {
	rnd := rand.New(1 << 23)
	ql := NewQlist()

	var (
		k float64
		v int
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		k = rnd.Float64()
		v = int(k * 10.0)
		if k > 0.5 {
			ql.Insert(k, v)
		} else if !ql.IsEmpty() {
			ql.PopMin()
		}
	}
}
