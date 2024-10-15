package main

type routeChecker struct {
	LimitDist   float64
	LimitTime   float64
	LimitNum    int
	LimitOrders int
	LimitWeight float64
	LimitVolume float64

	seenCoords map[LatLong]struct{}
}

func (rc *routeChecker) isFeasible2(route []int, uq map[LatLong]struct{}, dist, time float64) bool {
	if rc.LimitDist > 0 {
		if dist > rc.LimitDist {
			return false
		}
	}

	if rc.LimitTime > 0 {
		if time > rc.LimitTime {
			return false
		}
	}

	if rc.LimitNum > 0 {
		if len(uq) > rc.LimitNum {
			return false
		}
	}

	if rc.LimitOrders > 0 {
		if len(route)-1 > rc.LimitOrders {
			return false
		}
	}

	return true
}

func (rc *routeChecker) isFeasible(route []int, pts []Point, dm, tm [][]float64) bool {
	if rc.LimitDist > 0 {
		dist := 0.0
		for i := 1; i < len(route); i++ {
			dist += dm[pts[route[i-1]].Cid][pts[route[i]].Cid]
		}

		dist += dm[pts[route[len(route)-1]].Cid][0]

		if dist > rc.LimitDist {
			return false
		}
	}

	if rc.LimitTime > 0 {
		time := 0.0
		for i := 1; i < len(route); i++ {
			time += tm[pts[route[i-1]].Cid][pts[route[i]].Cid]
		}

		time += tm[pts[route[len(route)-1]].Cid][0]

		if time > rc.LimitTime {
			return false
		}
	}

	if rc.LimitNum > 0 {
		// count distance
		clear(rc.seenCoords)
		for i := 1; i < len(route); i++ {
			rc.seenCoords[pts[route[i]].Coord] = struct{}{}
		}

		if len(rc.seenCoords) > rc.LimitNum {
			return false
		}
	}

	if rc.LimitOrders > 0 {
		if len(route)-1 > rc.LimitOrders {
			return false
		}
	}

	return true
}

func isFeasible(routes [][]int, pts []Point, dm, tm [][]float64) bool {
	for _, r := range routes {
		if !rc.isFeasible(r, pts, dm, tm) {
			return false
		}
	}

	return true
}

func fitness(routes [][]int, pts []Point, dm [][]float64) float64 {
	f := 0.0
	for _, r := range routes {
		f += cost(r, pts, dm)
	}
	return f
}

func cost(route []int, pts []Point, dm [][]float64) float64 {
	dist := 0.0
	for i := 1; i < len(route); i++ {
		dist += dm[pts[route[i-1]].Cid][pts[route[i]].Cid]
	}

	if len(route) > 1 {
		dist += dm[pts[route[len(route)-1]].Cid][0]
	}

	return dist
}
