// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Understanding your workload is critically important in undserstanding
// if something can be made concurrent and how complex it is to perform.
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	numbers := generateList(1e7)

	fmt.Println(add(numbers))
	fmt.Println(addConcurrent(runtime.NumCPU(), numbers))
	fmt.Println(addConcurrent2(runtime.NumCPU(), numbers))
}

func generateList(totalNumbers int) []int {
	numbers := make([]int, totalNumbers)
	for i := 0; i < totalNumbers; i++ {
		numbers[i] = rand.Intn(totalNumbers)
	}
	return numbers
}

func add(numbers []int) int {
	var v int
	for _, n := range numbers {
		v += n
	}
	return v
}

func addConcurrent(goroutines int, numbers []int) int {
	var v int64
	totalNumbers := len(numbers)
	lastGoroutine := goroutines - 1
	stride := totalNumbers / goroutines

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			start := g * stride
			end := start + stride
			if g == lastGoroutine {
				end = totalNumbers
			}

			var lv int
			for _, n := range numbers[start:end] {
				lv += n
			}

			atomic.AddInt64(&v, int64(lv))
			wg.Done()
		}(g)
	}

	wg.Wait()

	return int(v)
}

func addConcurrent2(routines int, numbers []int) int {

	step := len(numbers) / routines
	ch := make(chan []int, routines)
	for i := 0; i < routines; i++ {
		start := i * step
		end := (i + 1) * step
		if end > len(numbers)-1 {
			end = len(numbers) - 1
		}

		ch <- numbers[start : end+1]
	}

	close(ch)

	res := make(chan int)

	var result int
	go func() {
		for sum := range res {
			result += sum
		}
	}()

	var wg sync.WaitGroup

	for nums := range ch {
		wg.Add(1)
		go func() {
			sum := 0
			for _, i := range nums {
				sum += i
			}
			res <- sum
			wg.Done()
		}()
	}

	wg.Wait()
	close(res)

	return result
}
