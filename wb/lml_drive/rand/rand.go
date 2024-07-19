package rand

import (
	"sync/atomic"
	"time"
)

const (
	a uint64 = 6364136223846793005 // Множитель
	c uint64 = 1442695040888963407 // Приращение
	m uint64 = 1 << 63             // Модуль, используется 2^63 для типа uint64
)

var (
	values []float64
	cap    int32
	idx    atomic.Int32
	seed   uint64
)

func Init(maxnum int) {
	seed = uint64(time.Now().UnixNano())

	buffer := make([]float64, 0, maxnum)

	for idx := 0; idx < maxnum; idx++ {
		buffer = append(buffer, nextFloat64())
	}

	values = buffer
	cap = int32(maxnum)
	idx = atomic.Int32{}
}

func next() uint64 {
	seed = (a*seed + c) % m
	return seed
}

func nextFloat64() float64 {
	return float64(next()) / float64(m)
}

func Float64() float64 {
	i := idx.Load()
	value := values[i]
	idx.Store((i + 1) % cap)

	return value
}

func Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to IntN")
	}
	return int(float64(n) * Float64())
}

// source from https://cs.opensource.google/go/go/+/refs/tags/go1.22.2:src/math/rand/v2/rand.go;drc=c29444ef39a44ad56ddf7b3d2aa8a51df1163e04;l=78
// func intn(n int) int {
// 	if n <= 0 {
// 		panic("invalid argument to IntN")
// 	}

// 	nuint := uint64(n)

// 	if nuint&(nuint-1) == 0 {
// 		return int(next() & (nuint - 1))
// 	}

// 	nxt := next()
// 	hi, lo := bits.Mul64(nxt, nuint)
// 	fmt.Println("hi, lo := bits.Mul64(next(), nuint)", nuint, nxt, hi, lo)
// 	if lo < nuint {
// 		thresh := -nuint % nuint
// 		fmt.Println("lo < nuint", nuint, hi, lo, thresh)
// 		for lo < thresh {
// 			hi, lo = bits.Mul64(next(), nuint)
// 		}
// 	}
// 	return int(hi)
// }
