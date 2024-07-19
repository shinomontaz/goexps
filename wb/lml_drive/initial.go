package main

import (
	"fmt"
	"sort"
)

func initialGreedy(pts []Point, dm [][]float64) [][]int {

	// pts[2525]
	idxs := make([]int, len(pts))
	for i := range pts {
		idxs[i] = i
	}

	bclst := ByClosest{list: idxs, pts: pts, tm: dm}
	sort.Sort(bclst)

	routes := make([][]int, 0) // маршрут = порядковые номера заказов в списке
	var curroute []int
	currShks := 0
	currDist := 0.0

	seenPoints := map[int]struct{}{0: struct{}{}}
	for _, p := range bclst.list[1:] {
		if _, ok := seenPoints[p]; ok {
			continue
		}

		curroute = []int{0, p}

		seenPoints[p] = struct{}{}
		currShks = pts[p].Shk
		currDist = dm[0][p]
		for {
			last := curroute[len(curroute)-1]
			closest := findClosestUnseen(last, bclst.list, seenPoints, dm)
			if closest == -1 {
				break // no points
			}
			_, maxs := fshk((currDist + dm[closest][0]) / 1000.0)
			if (currShks + pts[closest].Shk) > maxs {
				break // no points
			}
			curroute = append(curroute, closest)
			seenPoints[closest] = struct{}{}
			currShks += pts[closest].Shk
			currDist += dm[last][closest]
			if len(curroute) >= max_pts-1 {
				break
			}
		}

		routes = append(routes, curroute)
	}
	// проверка на дубли
	seen := map[int]struct{}{}
	for i, route := range routes {
		for _, p := range route {
			if _, ok := seen[p]; ok {
				fmt.Println(i, route, p)
				panic("duplicate!")
			}
			if p != 0 {
				seen[p] = struct{}{}
			}
		}
	}

	ptsOrderMap := map[int]struct{}{}
	for _, i := range idxs {
		ptsOrderMap[i] = struct{}{}
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

func findClosestUnseen(p int, pts []int, seen map[int]struct{}, tm [][]float64) int {
	res := -1
	min := -1.0
	for _, pp := range pts {
		if pp == p {
			continue
		}
		if _, ok := seen[pp]; ok {
			continue
		}

		if min == -1.0 {
			min = tm[p][pp] + tm[pp][0]
		}

		if min >= tm[p][pp]+tm[pp][0] {
			min = tm[p][pp] + tm[pp][0]
			res = pp
		}
	}

	return res
}
