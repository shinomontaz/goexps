package main

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/shinomontaz/goexps/Taboo/types"
)

func main() {

	tabooSize := 100
	iterations := 3 //10000

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
	current.edges = make([]Edge, 0, len(current.items)-1)
	current.orders = orders
	for idx := 1; idx < len(current.items); idx++ {
		current.edges = append(current.edges, Edge{current.items[idx-1], current.items[idx]})
	}

	current.Fitness = current.Distance()

	best := current
	tabuList := make([]Edge, 0, tabooSize)
	maxCandidates := 10

	for i := 0; i < iterations; i++ {
		candidates := make([]*Result, maxCandidates)
		for j := range candidates {
			candidates[j] = createCandidate(current, tabuList)
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Fitness < candidates[j].Fitness
		})

		for _, c := range candidates {
			fmt.Println("candidate", c.items, c.Fitness)
		}

		if candidates[0].Fitness < best.Fitness {
			best = candidates[0]

			fmt.Println(i, ":", best.items, best.Fitness)

			for _, edge := range best.edges {
				tabuList = append(tabuList, edge)
			}

			if len(tabuList) > tabooSize {
				tabuList = tabuList[(len(tabuList) - tabooSize):len(tabuList)]
			}
		}

		// отсортировать кандидатов по их фитнесу
		// если фитнесс лучшего кандидата лучше, чем текущий рекорд - обновить текущий рекорд кандидатом
		// добавить в tabooList список ребер кандидата
		// урезать голову tabooList так, чтобы общая длина списка была в рамках tabooSize

	}

	return best
}

func (r *Result) Clone() *Result {
	clone := &Result{
		items:   make([]int, 0, len(r.items)),
		Fitness: r.Fitness,
		edges:   make([]Edge, 0, len(r.edges)),
		orders:  r.orders,
	}

	for _, it := range r.items {
		clone.items = append(clone.items, it)
	}

	for _, e := range r.edges {
		clone.edges = append(clone.edges, e)
	}

	return clone
}

func (r *Result) isTabu(tabooList []Edge) bool {
	c2 := 0
	for idx, c := range r.items {
		if idx == len(r.items)-1 {
			c2 = r.items[0]
		} else {
			c2 = r.items[idx+1]
		}

		for _, tabooEdge := range tabooList {
			edge := Edge{
				start: c,
				end:   c2,
			}

			if tabooEdge == edge {
				return true
			}
		}
	}

	return false
}

func (r *Result) swap2opt(i, j int) {
	sol := r.items
	new_sol := make([]int, 0, len(sol))
	if i > j {
		i, j = j, i
	}

	new_sol = append(new_sol, sol[:i]...)
	for k := j - 1; k >= i; k-- {
		new_sol = append(new_sol, sol[k])
	}
	new_sol = append(new_sol, sol[j:]...)

	r.items = new_sol
	for idx := i + 1; idx < j; idx++ {
		r.edges[idx] = Edge{r.items[idx-1], r.items[idx]}
	}

}

func (r *Result) Mutate() *Result {
	copy := r.Clone()

	i := rand.Intn(len(r.items) - 2)
	j := i + rand.Intn(len(r.items)-i)

	if i == j {
		j = i + rand.Intn(len(r.items)-i)
	}

	copy.swap2opt(i, j)
	// update edges
	copy.Fitness = copy.Distance()

	return copy
}

func createCandidate(r *Result, tabooList []Edge) *Result {
	route := r
	i := 0
	for i = 0; i < 10000000; i++ {
		route = r.Mutate()
		if !route.isTabu(tabooList) {
			break
		}
	}
	fmt.Println(i)
	return route
}
