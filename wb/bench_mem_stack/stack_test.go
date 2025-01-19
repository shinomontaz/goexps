package main

import (
	"testing"
)

const size = 64 * 1024 //65536

// go test -bench=Insert -benchmem -benchtime=10s
func Benchmark_LargeSize_Stack_EqualOrLess65535(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// not escape to heap when size <= 65535
		dataLarge := make([]byte, size-1)
		_ = dataLarge
	}
}
func Benchmark_LargeSize_Heap_LargerThan65535(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// escape to heap when  size > 65535
		dataLarge := make([]byte, size+1)
		_ = dataLarge
	}
}
