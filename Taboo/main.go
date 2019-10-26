package main

import (
	"fmt"
	"math/rand"
	"sort"

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
	current.edges = make([]Edge, 0, len(current.items)-1)
	current.orders = orders
	for idx := 1; idx < len(current.items); idx++ {
		current.edges = append(current.edges, Edge{current.items[idx-1], current.items[idx]})
	}

	current.Fitness = current.Distance()

	fmt.Println("start: ", current.items, current.Fitness)

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

		panic("!")

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

func swap2opt(sol []int, i, j int) []int {
	new_sol := make([]int, 0, len(sol))
	if i > j {
		i, j = j, i
	}

	new_sol = append(new_sol, sol[:i]...)
	for k := j - 1; k >= i; k-- {
		new_sol = append(new_sol, sol[k])
	}
	new_sol = append(new_sol, sol[j:]...)
	return new_sol
}

func (r *Result) Mutate() *Result {
	copy := r.Clone()

	fmt.Println("before", copy.items)

	i := rand.Intn(len(r.items) - 2)
	j := i + rand.Intn(len(r.items)-i)

	if i == j {
		j = i + rand.Intn(len(r.items)-i)
	}

	copy.items = swap2opt(r.items, i, j)
	copy.Fitness = copy.Distance()

	fmt.Println("after", copy.items, i, j)

	return copy
}

func createCandidate(r *Result, tabooList []Edge) *Result {
	route := r
	for i := 0; i < 10000000; i++ {
		route := r.Mutate()
		if !route.isTabu(tabooList) {
			break
		}
	}
	return route
}
