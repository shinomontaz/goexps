package main

import "lml-drive/types"

func isFeasible(routes [][]int, tm [][]float64) int {
	max_over := 0
	for _, r := range routes {
		if len(r) > max_pts && max_over < len(r) {
			max_over = len(r)
		}
	}

	return max_over
}

// returns route distance, total overshk, total undershk, max overshk percent, max undershk
func fitness(routes [][]int, pts []types.Point, dm [][]float64) (float64, int, int, float64) {
	f := 0.0
	o := 0
	mx_o := 0.0
	u := 0
	rate := 0.0
	for _, r := range routes {
		ff, uu, oo, tshk := cost(r, pts, dm)
		f += ff
		o += oo
		u += uu
		if oo > 0 {
			rate = float64(oo) / float64(tshk-oo)
			if rate > mx_o {
				mx_o = rate
			}
		}
	}
	return f, o, u, mx_o
}

// dist, us, os, shks
func cost(route []int, pts []types.Point, dm [][]float64) (float64, int, int, int) {
	dist := 0.0
	total_shks := 0
	over_shk := 0
	undr_shk := 0
	for i := 1; i < len(route); i++ {
		dist += dm[route[i-1]][route[i]]
		total_shks += pts[route[i]].Shk
	}

	if len(route) > 1 {
		// добаляем время на возврат на склад
		dist += dm[route[len(route)-1]][0]
	}

	min_shk, max_shk := fshk(dist)

	if total_shks > max_shk {
		over_shk = total_shks - max_shk
	}
	if total_shks < min_shk {
		undr_shk = min_shk - total_shks
	}

	return dist, undr_shk, over_shk, total_shks
}
