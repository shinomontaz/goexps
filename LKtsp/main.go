package main

import (
	"fmt"
)

func main() {
	var distances [][]float64
	var ids []int
	m := &Matrix{d: distances, id: ids}

	m.solve()

	fmt.Println(m.tour)
}
