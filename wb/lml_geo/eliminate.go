package main

import (
	"fmt"
	"lml/rand"
)

func eliminate(routes [][]int, tm [][]float64) [][]int {

	// find n smallest routes to eliminate
	for i, r := range routes {
		if len(r) <= 1 {
			routes = append(routes[:i], routes[i+1:]...)
			fmt.Println("founded empty route on educate!")
			return routes
		}
	}

	lenTreshold := 1 + rand.Intn(5)

	suitable_routes := []int{}
	for rIdx, r := range routes {
		c, _ := cost(r, tm)
		if len(r) <= lenTreshold || c < 1.7*courier_shift {
			suitable_routes = append(suitable_routes, rIdx)
		}
	}

	if len(suitable_routes) == 0 {
		fmt.Println("no suitable_routes in elimintation")
		return routes
	}

	n := int(rand.Float64() * (float64(len(suitable_routes))))

	elimintating_routes := map[int]struct{}{}
	var random_suitable_route int

	for len(elimintating_routes) < n {
		// tournament selection
		random_suitable_route = 0
		for i := 0; i < 3; i++ {
			inst := rand.Intn(len(suitable_routes))
			if random_suitable_route == 0 || len(routes[suitable_routes[inst]]) < len(routes[suitable_routes[random_suitable_route]]) {
				random_suitable_route = inst
			}
		}

		elimintating_routes[random_suitable_route] = struct{}{}
	}

	pts := []int{}

	for i, r := range routes {
		if _, ok := elimintating_routes[i]; ok {
			for _, p := range r {
				if p != 0 {
					pts = append(pts, p)
				}
			}
		}
	}

	elimintated := 0
	for i := range elimintating_routes {
		key := i - elimintated
		routes = append(routes[:key], routes[key+1:]...)
		elimintated++
	}

	for _, p := range pts {
		// find proper route to add
		best_route, best_pos := findBestRouteToAdd(routes, p, tm)
		if best_pos < len(routes[best_route])-1 {
			routes[best_route] = append(routes[best_route][:best_pos+1], routes[best_route][best_pos:]...)
			routes[best_route][best_pos] = p
		} else {
			routes[best_route] = append(routes[best_route], p)
		}
	}

	ptsOrderMap := map[int]struct{}{}
	for _, p := range pts {
		ptsOrderMap[p] = struct{}{}
	}

	for _, route := range routes {
		for _, p := range route {
			if _, ok := ptsOrderMap[p]; ok {
				delete(ptsOrderMap, p)
			}
		}
	}

	if len(ptsOrderMap) > 0 {
		fmt.Println(ptsOrderMap)
		panic("init: not all points in solution!")
	}

	return routes
}

func eliminateRandom(routes [][]int, tm [][]float64) [][]int {
	for i, r := range routes {
		if len(r) <= 1 {
			routes = append(routes[:i], routes[i+1:]...)
			fmt.Println("founded empty route on eliminateRandom!")
			return routes
		}
	}

	random_suitable_route := rand.Intn(len(routes))

	pts := []int{}

	for _, p := range routes[random_suitable_route] {
		if p != 0 {
			pts = append(pts, p)
		}
	}

	routes = append(routes[:random_suitable_route], routes[random_suitable_route+1:]...)

	for _, p := range pts {
		// find proper route to add
		best_route, best_pos := findBestRouteToAdd(routes, p, tm)
		if best_pos < len(routes[best_route])-1 {
			routes[best_route] = append(routes[best_route][:best_pos+1], routes[best_route][best_pos:]...)
			routes[best_route][best_pos] = p
		} else {
			routes[best_route] = append(routes[best_route], p)
		}
	}

	ptsOrderMap := map[int]struct{}{}
	for _, p := range pts {
		ptsOrderMap[p] = struct{}{}
	}

	for _, route := range routes {
		for _, p := range route {
			if _, ok := ptsOrderMap[p]; ok {
				delete(ptsOrderMap, p)
			}
		}
	}

	if len(ptsOrderMap) > 0 {
		fmt.Println(ptsOrderMap)
		panic("eliminateRandom: not all points in solution!")
	}

	return routes
}

func findBestRouteToAdd(routes [][]int, so int, tm [][]float64) (route_id int, idx int) {
	newRoute := []int{}
	minOvertime := 0.0
	minLeg := 0.0
	var closestIdx int
	var minDist, currDist, currOvertime, currLeg float64
	for i, r := range routes { // найти наиболее подходящий маршрут для вставки
		closestIdx = 0
		minDist = tm[r[closestIdx]][so]

		for j, p := range r {
			currDist = tm[p][so]
			if currDist < minDist {
				closestIdx = j + 1
				minDist = currDist
			}
		}

		newRoute = newRoute[:0]
		for ii, pt := range r {
			newRoute = append(newRoute, pt)
			if ii == closestIdx {
				newRoute = append(newRoute, so)
			}
		}
		if closestIdx >= len(r)-1 {
			newRoute = append(newRoute, so)
		}

		currLeg, currOvertime = cost(newRoute, tm)
		if minLeg == 0 {
			minLeg = currLeg
			minOvertime = currOvertime
		}
		if ((currLeg < minLeg) && (currOvertime <= minOvertime)) || currOvertime < minOvertime {
			minLeg = currLeg
			route_id = i
			idx = closestIdx
		}
	}

	if idx == 0 {
		idx++
	}

	return route_id, idx
}
