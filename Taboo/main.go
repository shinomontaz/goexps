package main

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/shinomontaz/goexps/types"
)

func main() {

	tabooSize := 4
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

	best := current
	tabuList := make([]Edge, 0, tabooSize)
	maxCandidates := 20

	for i := 0; i < iterations; i++ {
		candidates := make([]*Result, maxCandidates)
		for j := range candidates {
			candidates[j] = createCandidate(best, tabuList)
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Fitness < candidates[j].Fitness
		})

		// for _, c := range candidates {
		// 	fmt.Println("candidate", c.items, c.Fitness)
		// }
		// fmt.Println(i, ":", candidates[0].Fitness, best.Fitness)

		// отсортировать кандидатов по их фитнесу
		// если фитнесс лучшего кандидата лучше, чем текущий рекорд - обновить текущий рекорд кандидатом
		// добавить в tabooList список ребер кандидата
		// урезать голову tabooList так, чтобы общая длина списка была в рамках tabooSize
		if candidates[0].Fitness < best.Fitness {
			best = candidates[0]

			fmt.Println(i, ":", best.Fitness)

			for _, edge := range best.edges {
				tabuList = append(tabuList, edge)
			}

			if len(tabuList) > tabooSize {
				tabuList = tabuList[(len(tabuList) - tabooSize):]
			}
		}
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
	//	fmt.Println("tabooList: ", tabooList)
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

			if tabooEdge.start == edge.start && tabooEdge.end == edge.end {
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

	//	fmt.Println(i, j, len(copy.items))
	// copy.items[i], copy.items[j] = copy.items[j], copy.items[i]

	// if i > 0 {
	// 	copy.edges[i-1] = Edge{copy.items[i-1], copy.items[i]}
	// }
	// copy.edges[i] = Edge{copy.items[i], copy.items[i+1]}

	// if j > 0 {
	// 	copy.edges[j-1] = Edge{copy.items[j-1], copy.items[j]}
	// }
	// if j < len(copy.items)-1 {
	// 	copy.edges[j] = Edge{copy.items[j], copy.items[j+1]}
	// }

	copy.swap2opt(i, j)
	// update edges
	copy.Fitness = copy.Distance()

	return copy
}

func createCandidate(r *Result, tabooList []Edge) *Result {
	route := r.Clone()
	i := 0
	for i = 0; i < 1000000; i++ {
		route = route.Mutate()
		if !route.isTabu(tabooList) {
			//			fmt.Println("route: ", route)
			break
		}
		// if i == 100 {
		// 	fmt.Println("tabooList", tabooList, "route.items", route.items, route.edges)
		// 	panic("!taboo")
		// }
	}

	//	fmt.Println(i)
	return route
}
