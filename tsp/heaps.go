package tsp

import (
	"fmt"
	"math"
)

func solve(
	dm [][]int64,
	first,
	last int,
) ([]int, error) {
	if first > len(dm)-1 || last > len(dm)-1 {
		return nil, fmt.Errorf("border points can't be more than len(dm)")
	}
	if len(dm) < 2 {
		return nil, fmt.Errorf("distance matrix size must be greater or equal to 2")
	}

	var res []int

	offset := 0
	offset2 := 0
	if last > 0 && last != first {
		offset = 1
	}
	if last == first { // кольцевой маршрут длинее на 1
		offset2 = 1
	}
	upts := make([]int, len(dm)-1-offset) // - 1 из-за явной начальной точки, - offset если не цикличный маршрут
	res = make([]int, len(dm)+offset2)
	j := 0
	for i := range dm {
		if i == first || i == last {
			continue
		}
		upts[j] = i
		j++
	}

	var min int64
	var curCost int64
	currvar := make([]int, len(upts))

	min = math.MaxInt64

	for p := make([]int, len(upts)); p[0] < len(p); nextPerm(p) {
		getPerm(upts, p, currvar)

		curCost = dm[first][currvar[0]]
		if curCost >= min {
			continue
		}
		for i := 1; i < len(currvar); i++ {
			curCost += dm[currvar[i-1]][currvar[i]]
			if curCost >= min {
				continue
			}
		}
		curCost += dm[currvar[len(currvar)-1]][last]

		//		curCost = findCost(currvar, dm, first, last)
		if min > curCost {
			min = curCost
			// form result
			res[0] = first
			for i := range currvar {
				res[i+1] = currvar[i]
			}
			res[len(res)-1] = last
		}
	}

	return res, nil
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int, result []int) []int {
	copy(result, orig)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

func findCost(route []int, costMatrix [][]int64, first, last int) (cost int64) {
	cost = costMatrix[first][route[0]]
	for i := 1; i < len(route); i++ {
		cost += costMatrix[route[i-1]][route[i]]
	}
	cost += costMatrix[route[len(route)-1]][last]

	return
}
