package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	tabooSize     = 4
	iterations    = 10000
	maxCandidates = 20
)

func main() {
	route, fitness := search(ords, dm)

	fmt.Println(route, fitness)
}

type Edge struct {
	start int
	end   int
}

func findCost(route []int, dm [][]float64) float64 {
	var curCost float64
	for i := 1; i < len(route); i++ {
		curCost = curCost + dm[route[i-1]][route[i]]
	}
	return curCost
}

func search(upts []LatLng, dm [][]float64) ([]int, float64) {
	indices := make([]int, len(upts))
	// for idx := 1; idx < len(current.items); idx++ {
	// 	current.edges = append(current.edges, Edge{current.items[idx-1], current.items[idx]})
	// }
	for i := range upts {
		indices[i] = i
	}

	best := indices
	bestDist := findCost(best, dm)

	tabuList := make([]Edge, 0, tabooSize)

	candidates := make([][]int, maxCandidates)
	candDistances := make([]float64, maxCandidates)
	candEdges := make([][]Edge, maxCandidates)

	for i := 0; i < iterations; i++ {
		for c := 0; c < maxCandidates; c++ {
			candidates[c], candEdges[c] = createCandidate(best, tabuList, iterations)
			candDistances[c] = findCost(candidates[c], dm)
		}

		// поиск лучшего решения без сортировки
		bestId := -1
		for id := range candidates {
			if candDistances[id] < bestDist {
				bestId = id
				bestDist = candDistances[bestId]
			}
		}

		if bestId >= 0 {
			best = candidates[bestId]
			tabuList = append(candEdges[bestId], tabuList...)
			if len(tabuList) > tabooSize {
				tabuList = tabuList[(len(tabuList) - tabooSize):]
			}
		}

	}

	return best, bestDist
}

func createCandidate(r []int, tabooList []Edge, iters int) ([]int, []Edge) {
	var lenR = len(r)
	route := make([]int, lenR)
	edges := make([]Edge, lenR-1)
	reps := 0

NextCand:
	for reps = 0; reps < iters; reps++ {
		copy(route, r)
		i := 1 + rs.Intn(lenR-2)
		j := i + rs.Intn(lenR-i)
		if i == j {
			continue
		}
		// делаем swap2opt
		if i > j {
			i, j = j, i
		}
		if j-i == 1 {
			continue
		}
		for k := 0; k < (j-i)/2; k++ {
			route[i+k], route[j-k] = route[j-k], route[i+k]
		}

		for i = 0; i < lenR-1; i++ {
			edges[i] = Edge{route[i], route[i+1]}
		}

		// проверить, что полученное решение не в taboo
		for _, edge := range edges {
			for _, taboo := range tabooList {
				if taboo.start == edge.start && taboo.end == edge.end {
					continue NextCand
				}
			}
		}
		return route, edges
	}
	logrus.Warn("creating candidate end by repeats count! return last..")
	return route, edges
}
