package main

import (
	"math/rand"
	"testing"
	"time"

	oddg "github.com/oddg/hungarian-algorithm"
)

var rs *rand.Rand

func init() {
	rs = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

const MAX = 1000

func GenerateMatrix(n int) [][]int {

	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		for j := 0; j < n; j++ {
			res[i][j] = rs.Intn(MAX)
		}
	}

	return res
}

// >go test -bench=Algo -v -run=^$ bench_test.go

func BenchmarkAlgo1(b *testing.B) {
	Num := 100

	b.StopTimer()
	matrix := GenerateMatrix(Num)
	b.Log(Num)
	b.StartTimer()
	oddg.Solve(matrix)
}

func BenchmarkAlgo2(b *testing.B) {
	Num := 100

	b.StopTimer()
	matrix := GenerateMatrix(Num)
	b.Log(Num)
	b.StartTimer()
	myAnswer := algorithm(matrix)
}
