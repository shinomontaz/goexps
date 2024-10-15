package main

import (
	"fmt"
	"sort"
)

func initialGreedy(pts []Point, dm [][]float64) [][]int {
	bclst := ByClosest{pts: pts, tm: dm}
	sort.Sort(bclst)

	routes := make([][]int, 0) // маршрут = порядковые номера заказов в списке
	var curroute []int
	currDist := 0.0

	uq := map[LatLong]struct{}{}
	seenPoints := map[int]struct{}{0: struct{}{}}
	for i, p := range pts {
		if _, ok := seenPoints[i]; ok {
			continue
		}
		clear(uq)
		curroute = []int{0, i}

		seenPoints[i] = struct{}{}
		currDist = dm[0][p.Cid]
		uq[p.Coord] = struct{}{}
		for {
			last := curroute[len(curroute)-1]
			closest := findClosestUnseen(last, pts, seenPoints, dm)
			if closest == -1 {
				break // no points
			}

			curroute = append(curroute, closest)
			seenPoints[closest] = struct{}{}
			currDist += dm[pts[last].Cid][pts[closest].Cid]
			uq[pts[closest].Coord] = struct{}{}

			if !rc.isFeasible2(curroute, uq, currDist+dm[pts[closest].Cid][0], 0) {
				curroute = curroute[:len(curroute)-1]
				delete(uq, pts[closest].Coord)
				delete(seenPoints, closest)
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
	for i := range pts {
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

func findClosestUnseen(p int, pts []Point, seen map[int]struct{}, tm [][]float64) int {
	res := -1
	min := -1.0
	for i, pp := range pts {
		if i == p {
			continue
		}
		if _, ok := seen[i]; ok {
			continue
		}

		if min == -1.0 {
			min = tm[pts[p].Cid][pp.Cid] + tm[pp.Cid][0]
		}

		if min >= tm[pts[p].Cid][pp.Cid]+tm[pp.Cid][0] {
			min = tm[pts[p].Cid][pp.Cid] + tm[pp.Cid][0]
			res = i
		}
	}

	return res
}
