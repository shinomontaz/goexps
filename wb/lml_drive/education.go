package main

import (
	"lml-drive/rand"
)

func educate(routes [][]int, pts []Point, dm [][]float64) [][]int {
	for i, r := range routes {
		if len(r) <= 1 {
			routes = append(routes[:i], routes[i+1:]...)
			//			fmt.Println("founded empty route on educate!")
			return routes
		}
	}

	suitable_routes := []int{}
	for rIdx, r := range routes {
		dist, os, _, shks := cost(r, pts, dm)
		mdist := fdist(dist)
		shk_min, _ := fshk(dist)
		if os == 0 && ((dist <= mdist/2) || (shks <= shk_min)) {
			suitable_routes = append(suitable_routes, rIdx)
		}
	}

	if len(suitable_routes) == 0 {
		return routes
	}

	// tournament selection
	var random_suitable_route int
	for i := 0; i < 5; i++ {
		inst := rand.Intn(len(suitable_routes))
		if random_suitable_route == 0 || len(routes[suitable_routes[inst]]) < len(routes[suitable_routes[random_suitable_route]]) {
			random_suitable_route = inst
		}
	}

	rIdx := suitable_routes[random_suitable_route]

	newRoute := []int{}
	for _, it := range routes[rIdx] { // move items to other routes
		if it == 0 {
			continue
		}
		// найти для it наиболее подходящий маршрут. Такой, где а) менее всего нарушается over_shk и б) наименьшая общая длина получается
		needId := rIdx
		minDist := 0.0
		minOvshk := -1
		for id, r := range routes {
			if id == rIdx {
				continue
			}
			newRoute = []int{}
			newRoute = append(newRoute, r...)
			newRoute = append(newRoute, it)
			dist, ov, _, _ := cost(newRoute, pts, dm)
			if minOvshk == -1 || ov < minOvshk {
				minOvshk = ov
				minDist = dist
				needId = id
			} else if ov == minOvshk {
				if dist < minDist {
					minOvshk = ov
					minDist = dist
					needId = id
				}
			}
		}

		routes[needId] = append(routes[needId], it)
	}

	i := 0
	for idx, r := range routes {
		if idx != rIdx {
			routes[i] = make([]int, len(r))
			copy(routes[i], r)
			i++
		}
	}
	routes = routes[:i]

	return routes
}
