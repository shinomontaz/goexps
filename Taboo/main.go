package main

import (
	"fmt"

	"github.com/shinomontaz/goexps/Taboo/types"
)

func main() {

	tabooSize := 10
	iterations := 10000

	result := search(orders, tabooSize, iterations)

	fmt.Println(result)
}

type Edge struct {
	start int
	end   int
}

type Result struct {
	items   []int
	Fitness float64
	edges   []Edge
}

func search(orders []*types.Order, tabooSize int, iterations int) *Result {
	res := &Result{}

	return res
}
