package main

import (
	"fmt"
	"testing"
	"time"

	oddg "github.com/oddg/hungarian-algorithm"
)

func TestCompare(t *testing.T) {
	matrix := GenerateMatrix(1000)
	//	fmt.Println(matrix)
	var start time.Time
	var elapsed time.Duration
	start = time.Now()
	//	oddgAnswer, _ := oddg.Solve(matrix)
	oddg.Solve(matrix)
	elapsed = time.Since(start)
	//	fmt.Println("oddgAnswer", oddgAnswer)
	fmt.Println("oddg takes: ", elapsed)

	start = time.Now()
	//	myAnswer := algorithm(matrix)
	algorithm(matrix)
	elapsed = time.Since(start)
	//	fmt.Println("myAnswer", myAnswer)
	fmt.Println("my takes: ", elapsed)
}
