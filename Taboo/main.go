package main

import (
	"fmt"
	"math/rand"

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
	orders  []*types.Order
}

func (r *Result) Distance() float64 {
	res := 0.0

	for _, e := range r.edges {
		res += r.orders[e.start].Coord.Distance(r.orders[e.end].Coord)
	}

	return res
}

func search(orders []*types.Order, tabooSize int, iterations int) *Result {
	current := &Result{}
	current.items = rand.Perm(len(orders))
	current.edges = make([]Edge, 0, len(res.items)-1)
	current.orders = orders
	for idx, i := range res.items {
		if idx == 0 {
			continue
		}
		current.edges = append(res.edges, Edge{i - 1, i})
	}
	current.Fitness = res.Distance()

	best := current
	tabuList := make([]Edge, 0, tabooSize)
	maxCandidates := 10

	for i := 0; i < iterations; i++ {
		candidates := make([]Result, maxCandidates)
		for j := range candidates {
			candidates[j] = createCandidate(current, tabuList, orders)
		}

		// отсортировать кандидатов по их фитнесу
		// если фитнесс лучшего кандидата лучше, чем текущий рекорд - обновить текущий рекорд кандидатом
		// добавить в tabooList список ребер кандидата
		// урезать голову tabooList так, чтобы общая длина списка была в рамках tabooSize

	}

	return best
}

func createCandidate(r *Result, tabooList []Edge, orders []*types.Order) Result {
	for i := 0; i < 10000000; i++ {
		route := r.Mutate()
		if !route.isTabu(tabooList) {
			break
		}
	}
	return *route
}
