package main

func isFeasible(routes [][]int) bool {
	for _, r := range routes {
		if len(r) > max_pts {
			return false
		}
	}

	return true
}

func fitness(routes [][]int, pts []Point, dm [][]float64) (float64, int, int) {
	f := 0.0
	o := 0
	u := 0
	for _, r := range routes {
		ff, oo, uu, _ := cost(r, pts, dm)
		f += ff
		o += oo
		u += uu
	}
	return f, o, u
}

func cost(route []int, pts []Point, dm [][]float64) (float64, int, int, int) {
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

	min_shk, max_shk := fshk(dist / 1000.0)

	if total_shks > max_shk {
		over_shk = total_shks - max_shk
	}
	if total_shks < min_shk {
		undr_shk = min_shk - total_shks
	}

	return dist, undr_shk, over_shk, total_shks
}
