package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func printAllocInfo(message string) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Println(message,
		"Heap Alloc:", memStats.HeapAlloc,
		"Heap objects:", memStats.HeapObjects)
}

func createMapper(n int) []map[int]int {
	m := make([]map[int]int, n)
	for i := range m {
		m[i] = make(map[int]int)
	}

	return m
}

func createSlicer(n int) [][]int {
	s := make([][]int, n)
	for i := range s {
		s[i] = []int{0}
	}

	return s
}

func main() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	mm := createMapper(100)
	ss := createSlicer(100)
	for i := range mm {
		randLen := rnd.Intn(10)
		for j := range randLen {
			mm[i][j] = rnd.Intn(10)
		}
	}

	for i := range ss {
		randLen := rnd.Intn(10)
		for j := range randLen {
			ss[i] = append(ss[i], j)
		}
	}

	min_val := 10
	undercount := 0
	printAllocInfo("before map loop")
	for i := 0; i < 1000000; i++ {
		undercount = 0
		for _, m := range mm {
			sum := 0
			for _, v := range m {
				if v > 2 {
					sum += v
				}
			}

			if sum < min_val {
				undercount += min_val - sum
			}
		}
	}
	printAllocInfo("after map loop")
	fmt.Println(undercount)

	sum := 0
	printAllocInfo("before range loop")
	for i := 0; i < 1000000; i++ {
		undercount = 0
		for _, s := range ss {
			sum = 0
			for _, v := range s {
				if v != 0 {
					sum += 1
				}
			}

			if sum < min_val {
				undercount += min_val - sum
			}
		}
	}
	printAllocInfo("after range loop")
	fmt.Println(undercount)
}
