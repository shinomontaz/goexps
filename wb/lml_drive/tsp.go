package main

import "lml-drive/types"

func swapRevert(points []int, i, j int) []int {
	for k := 0; k < (j-i+1)/2; k++ {
		points[i+k], points[j-k] = points[j-k], points[i+k]
	}

	return points
}

func findCost(ids []int, costMatrix [][]float64) (cost float64) {
	for i := 1; i < len(ids); i++ {
		cost += costMatrix[ids[i-1]][ids[i]]
	}

	return
}

func swap2opt(upts []types.Point, dm [][]float64) ([]int, float64) {
	route := make([]int, len(upts))
	for i := range upts {
		route[i] = i
	}
	dist := findCost(route, dm)

	is_found := true
	n := len(upts)
	for is_found {
		is_found = false
		for j := 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				new_route := swapRevert(route, j, k)
				new_distance := findCost(new_route, dm)
				if new_distance < dist {
					dist = new_distance
					route = new_route
					is_found = true

					break
				}
			}
		}
	}

	return route, dist
}
