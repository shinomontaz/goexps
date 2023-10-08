package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	oddg "github.com/oddg/hungarian-algorithm"
)

const MAX = 1000

var rs *rand.Rand

func init() {
	rs = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

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

func TestCompare(t *testing.T) {
	matrix := GenerateMatrix(100)
	fmt.Println(matrix)
	var start time.Time
	var elapsed time.Duration
	start = time.Now()
	oddgAnswer, _ := oddg.Solve(matrix)
	elapsed = time.Since(start)
	fmt.Println("oddgAnswer", oddgAnswer)
	fmt.Println("oddg takes: ", elapsed)

	start = time.Now()
	myAnswer := algorithm(matrix)
	elapsed = time.Since(start)
	fmt.Println("myAnswer", myAnswer)
	fmt.Println("my takes: ", elapsed)
}
