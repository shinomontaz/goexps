package main

/*
import "lml/rand"

func lns(routes [][]int, dm, tm [][]float64) [][]int {
	// 1. choose routes to stay
	dice := rand.Float64()
	n := int(float64(len(routes)) * dice / 2.0)
	if dice < 0.5 {
		idx_routes_to_stay = best_choose(n, routes, dm, tm)
	} else {
		idx_routes_to_stay = random_choose(n, routes)
	}

	newRoutes = newRoutes[:0]

	clear(ptsToLns)
	for i, r := range routes {
		if _, ok := idx_routes_to_stay[i]; !ok {
			for _, p := range r {
				ptsToLns[p] = struct{}{}
			}
		} else {
			newRoutes = append(newRoutes, make([]int, len(r)))
			copy(newRoutes[len(newRoutes)-1], r)

			newSupplMap = append(newSupplMap, map[int]int{})
			for k, v := range supplMap[i] {
				if v > 0 {
					newSupplMap[len(newSupplMap)-1][k] = v
				}
			}
		}
	}

	// construct from ptsToLns something new
	ptsOrder := make([]int, 0, len(ptsToLns))
	ptsOrder = append(ptsOrder, 0)

	for id := range ptsToLns {
		if ptsMap[id].Type == TYPE_CLIENT { // поставщика все равно добавим куда надо
			ptsOrder = append(ptsOrder, id)
		}
	}
}
*/
