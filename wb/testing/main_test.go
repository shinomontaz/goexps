package main

import (
	"testing"
)

var num = 1000

func BenchmarkPrimeNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		primeNumbers(num)
	}
}
