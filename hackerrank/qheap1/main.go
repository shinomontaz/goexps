package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var heap []int

func main() {
	//Enter your code here. Read input from STDIN. Print output to STDOUT

	var cmd string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cmd = scanner.Text()

	queries, _ := strconv.Atoi(cmd)

	for i := 0; i < queries; i++ { // endless loop to get input data
		scanner.Scan()
		cmd = scanner.Text()
		process(cmd)
	}
}

func process(cmd string) {
	args := strings.Split(cmd, " ")
	op, _ := strconv.Atoi(args[0])

	switch op {
	case 1: // append
		el, _ := strconv.Atoi(args[1])
		add(el)
	case 2: // delete
		idx, _ := strconv.Atoi(args[1])
		del(idx)
	case 3: // print
		prnt()
	}
}

func add(e int) {
	heap = append(heap, e) // add to tail
	// sift_up
	sift_up(len(heap) - 1)
}

func prnt() {
	fmt.Printf("%d\n", heap[0])
}

func del(e int) {
	// find e in heap and delete it
	var id int
	for idx, el := range heap {
		if el == e {
			id = idx
			break
		}
	}

	heap[id], heap[len(heap)-1] = heap[len(heap)-1], heap[id]
	heap = heap[:len(heap)-1]

	sift_down(id)
}

func sift_up(idx int) {
	if idx <= 0 {
		return
	}

	parent := int((idx - 1) / 2)
	if heap[parent] > heap[idx] {
		heap[parent], heap[idx] = heap[idx], heap[parent]
		sift_up(parent)
	}
}

func sift_down(idx int) {
	n := len(heap) - 1
	if idx >= n-1 {
		return
	}

	//find 2 childs:
	min := heap[idx]
	new_idx := idx
	if 2*idx+1 < n && min > heap[2*idx+1] {
		min = heap[2*idx+1]
		new_idx = 2*idx + 1
	}
	if 2*idx+2 < n && min > heap[2*idx+2] {
		min = heap[2*idx+2]
		new_idx = 2*idx + 2
	}

	if new_idx == idx {
		return
	}

	heap[idx], heap[new_idx] = heap[new_idx], heap[idx]

	sift_down(new_idx)
}
