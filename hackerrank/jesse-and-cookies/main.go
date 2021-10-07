package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'cookies' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts following parameters:
 *  1. INTEGER k
 *  2. INTEGER_ARRAY A
 */

func cookies(k int32, A []int32) int32 {
	// Write your code here
	// make heap from A
	heap := make([]int64, 0)
	for _, e := range A {
		add(int64(e), &heap)
	}

	var steps int32

	for len(heap) > 1 && heap[0] <= int64(k) {
		min1 := del(&heap)
		min2 := del(&heap)

		var newcookie int64
		if min1 < min2 {
			newcookie = int64(min1 + min2*2)
		} else {
			newcookie = int64(min2 + min1*2)
		}

		add(newcookie, &heap)
		steps++
	}

	if heap[0] >= int64(k) {
		return steps
	}

	//	fmt.Printf("%v, %v\n", heap[0], heap)
	return -1
}

func add(e int64, heap *[]int64) {
	*heap = append(*heap, e) // add to tail
	sift_up(len(*heap)-1, heap)
}

func del(heap *[]int64) int64 { // always delete root
	res := (*heap)[0]
	(*heap)[0], (*heap)[len(*heap)-1] = (*heap)[len(*heap)-1], (*heap)[0]
	*heap = (*heap)[:len(*heap)-1]
	sift_down(0, heap)

	return res
}

func sift_up(idx int, heap *[]int64) {
	if idx <= 0 {
		return
	}

	parent := int((idx - 1) / 2)
	if (*heap)[parent] > (*heap)[idx] {
		(*heap)[parent], (*heap)[idx] = (*heap)[idx], (*heap)[parent]
		sift_up(parent, heap)
	}
}

func sift_down(idx int, heap *[]int64) {
	n := int(len(*heap) - 1)
	if idx >= n-1 {
		return
	}

	//find 2 childs:
	min := (*heap)[idx]
	new_idx := idx
	if 2*idx+1 < n && min > (*heap)[2*idx+1] {
		min = (*heap)[2*idx+1]
		new_idx = 2*idx + 1
	}
	if 2*idx+2 < n && min > (*heap)[2*idx+2] {
		min = (*heap)[2*idx+2]
		new_idx = 2*idx + 2
	}

	if new_idx == idx {
		return
	}

	(*heap)[idx], (*heap)[new_idx] = (*heap)[new_idx], (*heap)[idx]

	sift_down(new_idx, heap)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)
	//    stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	stdout, err := os.Create("out.txt")

	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	firstMultipleInput := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	nTemp, err := strconv.ParseInt(firstMultipleInput[0], 10, 64)
	checkError(err)
	n := int32(nTemp)

	kTemp, err := strconv.ParseInt(firstMultipleInput[1], 10, 64)
	checkError(err)
	k := int32(kTemp)

	ATemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	var A []int32

	for i := 0; i < int(n); i++ {
		AItemTemp, err := strconv.ParseInt(ATemp[i], 10, 64)
		checkError(err)
		AItem := int32(AItemTemp)
		A = append(A, AItem)
	}

	result := cookies(k, A)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
