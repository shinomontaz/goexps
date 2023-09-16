package main

import "fmt"

func main() {

	arr := []int{1, 3, 5, 7, 8, 9, 10, 14, 15, 25, 27, 29}
	x := 16 // search min i | a[i] >= x

	l := -1
	r := len(arr)
	fmt.Println("r", r)

	var m int
	for r-l > 1 {
		m = int((r + l) / 2)

		if arr[m] < x {
			l = m
		} else {
			r = m
		}
	}

	if r == len(arr) {
		// not finded
	} else {
		fmt.Println("finded: ", arr[r], r)
	}
}
