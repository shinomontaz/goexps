package main

import (
	"fmt"
	"lml-drive/types"

	_ "net/http/pprof"
)

func solve(pts []types.Point, last_route_id int) (gj DataRouted, table [][]string, overshk int, undershk int) {
	// матрицы расстояний
	l := len(pts)
	dm := make([][]float64, l)
	tm := make([][]float64, l)

	for i := 0; i < l; i++ {
		dm[i] = make([]float64, l)
		tm[i] = make([]float64, l)
	}

	max_len := 500
	chunks := int(float64(l) / float64(max_len))
	if l%max_len > 0 {
		chunks++
	}

	var src, dst []int
	for j := 0; j < chunks; j++ {
		src = make([]int, 0, l)
		for id := j * max_len; id < l && id < (j+1)*max_len; id++ {
			src = append(src, id)
		}
		fmt.Println("start receiving ", j, " out of ", chunks)
		for jj := 0; jj < chunks; jj++ {
			fmt.Println("start receiving ", jj, " out of ", chunks, " for row ", j)

			dst = make([]int, 0, l)
			for id := jj * max_len; id < l && id < (jj+1)*max_len; id++ {
				dst = append(dst, id)
			}
			ttm, tdm := GetMatrices2(pts, src, dst) // times and distances
			for iii, row := range tdm {
				for jjj, p := range row {
					dm[iii+j*max_len][jjj+jj*max_len] = p
					tm[iii+j*max_len][jjj+jj*max_len] = ttm[iii][jjj]
				}
			}
		}
	}

	initialRoutes := initialGreedy(pts, dm, tm)

	currSolution := make([][]int, len(initialRoutes))
	for i, r := range initialRoutes {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}
	currSolution = sa(currSolution, pts, dm, tm, 0.9999)

	for _, r := range currSolution {
		dist, ushk, oshk, shks := cost(r, pts, dm)
		fmt.Println("dist: ", dist, "undershk: ", ushk, "overshk: ", oshk, shks)
	}

	newSolution := educate(currSolution, pts, dm, 0.8)
	if len(newSolution) < len(currSolution) {
		fmt.Println("!!!", len(newSolution), len(currSolution))
	}

	currSolution = currSolution[:len(newSolution)]
	copySlice(currSolution, newSolution)

	_, overshk, undershk, _ = fitness(currSolution, pts, dm)

	return prepareGeoJson(currSolution, pts, dm, tm), prepareFile(currSolution, pts, tm, last_route_id), overshk, undershk

	//	fmt.Println(len(currSolution), oldFitness, oldOvershk, oldUndershk)
}
