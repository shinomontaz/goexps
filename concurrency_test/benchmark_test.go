package main

import (
	"runtime"
	"testing"
)

var numbers []int

func init() {
	numbers = generateList(1e7)
}

func BenchmarkSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add(numbers)
	}
}

func BenchmarkConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addConcurrent(runtime.NumCPU(), numbers)
	}
}

func BenchmarkConcurrent2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addConcurrent2(runtime.NumCPU(), numbers)
	}
}
