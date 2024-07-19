package main

import (
	"fmt"
	"lml/rand"
	"sort"
)

func initialFlower(pts []LatLong, tm, dm [][]float64) [][]int {

	// sort by polar around 0 point

	idxs := make([]int, len(pts))
	for i := range pts {
		idxs[i] = i
	}

	bplr := ByPolar{list: idxs, pts: pts}
	sort.Sort(bplr)

	idxs = bplr.list

	currIdx := rand.Intn(len(pts))
	rndDirection := 1 - rand.Intn(2)*2

	routes := make([][]int, 0) // маршрут = порядковые номера заказов в списке

	curroute := []int{0}
	lenOrds := len(pts)
	currTime := 0.0

	var (
		dt       float64
		addition float64
	)
	for seen := 0; seen < lenOrds; seen++ {
		ptId := idxs[currIdx]
		dt = tm[curroute[len(curroute)-1]][currIdx]

		addition = 0
		if len(curroute) > 1 {
			addition = tm[currIdx][0]
		}

		if currTime+dt+box_time+addition > courier_shift {
			routes = append(routes, curroute)
			curroute = []int{0} // начало на складе WB

			currTime = 0.0
			dt = tm[0][currIdx]
		}

		curroute = append(curroute, ptId)
		currTime += dt + box_time

		currIdx += rndDirection
		if currIdx >= lenOrds {
			currIdx -= lenOrds
		}
		if currIdx < 0 {
			currIdx = lenOrds + currIdx
		}

	}

	routes = append(routes, curroute)

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

func initialGreedy(pts []LatLong, tm [][]float64) [][]int {

	// pts[2525]
	idxs := make([]int, len(pts))
	for i := range pts {
		idxs[i] = i
	}

	bclst := ByClosest{list: idxs, pts: pts, tm: tm}
	sort.Sort(bclst)

	routes := make([][]int, 0) // маршрут = порядковые номера заказов в списке
	var curroute []int
	currTime := 0.0
	dt := 0.0
	seenPoints := map[int]struct{}{0: struct{}{}}
	for _, p := range bclst.list[1:] {
		if _, ok := seenPoints[p]; ok {
			continue
		}

		curroute = []int{0, p}

		seenPoints[p] = struct{}{}
		currTime = box_time + tm[0][p]
		for {
			last := curroute[len(curroute)-1]
			closest := findClosestUnseen(last, bclst.list, seenPoints, tm)
			if closest == -1 {
				break // no points
			}
			dt = tm[last][closest]
			if dt > 0 {
				dt += box_time
			}
			curroute = append(curroute, closest)
			currTime += dt
			seenPoints[closest] = struct{}{}

			if currTime+tm[closest][0] > courier_shift {
				break
			}

			// if currTime+dt+tm[closest][0] <= courier_shift {
			// 	curroute = append(curroute, closest)
			// 	currTime += dt + box_time
			// 	seenPoints[closest] = struct{}{}
			// } else {
			// 	break
			// }
		}

		routes = append(routes, curroute)
		// find other closest points to the end
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

func splitMaxRoute(routes [][]int, num int, tm [][]float64) [][]int {
	// find num routes with maximum overtime
	bovtm := ByOvertime{
		routes: routes,
		tm:     tm,
		list:   make([]int, len(routes)),
	}

	for i := range routes {
		bovtm.list[i] = i
	}

	sort.Sort(bovtm)

	splitting_routes := bovtm.list[:num]
	var (
		pts      []int
		currTime float64
		idx      int
		dt       float64
	)

	for _, i := range splitting_routes {
		pts = []int{0}
		currTime = 0
		idx = 0
		for j, p := range routes[i][1:] {
			dt = tm[routes[i][j]][p]
			if dt > 0 {
				dt += box_time
			}
			currTime += dt

			if currTime > courier_shift {
				idx = j
				pts = append(pts, p)
			}
		}

		if idx == 0 {
			idx = len(routes[i]) - 1
			pts = append(pts, routes[i][idx])
		}

		routes[i] = routes[i][:idx]
		routes = append(routes, pts)
	}

	return routes
}

func initialRandomGreedy(pts []LatLong, tm [][]float64) [][]int {
	idxs := make([]int, len(pts))
	for i := range pts {
		idxs[i] = i
	}

	for i := len(idxs) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		idxs[i], idxs[j] = idxs[j], idxs[i]
	}

	routes := make([][]int, 0) // маршрут = порядковые номера заказов в списке
	var curroute []int
	currTime := 0.0
	dt := 0.0
	seenPoints := map[int]struct{}{0: struct{}{}}
	for _, p := range idxs {
		if _, ok := seenPoints[p]; ok {
			continue
		}

		curroute = []int{0, p}

		seenPoints[p] = struct{}{}
		currTime = box_time + tm[0][p]
		for {
			last := curroute[len(curroute)-1]
			closest := findClosestUnseen(last, idxs, seenPoints, tm)
			if closest == -1 {
				break // no points
			}
			dt = tm[p][closest]
			if dt > 0 {
				dt += box_time
			}

			curroute = append(curroute, closest)
			currTime += dt
			seenPoints[closest] = struct{}{}
			if currTime > courier_shift {
				break
			}
			// if currTime+dt < courier_shift {
			// 	curroute = append(curroute, closest)
			// 	currTime += dt
			// 	seenPoints[closest] = struct{}{}
			// } else {
			// 	break
			// }

		}

		routes = append(routes, curroute)
		// find other closest points to the end
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
