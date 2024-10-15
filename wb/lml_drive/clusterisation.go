package main

import (
	"fmt"
	"lml-drive/types"
)

func clusterize(start_num int, pts []types.Point) [][]types.Point {
	// Partition the data points into 16 clusters
	km := NewKmeans()
	res := make([][]types.Point, 0)
	skip := false
	sum_shk := 0
	for num := start_num; num > 0; num-- {
		clusters, err := km.Do(pts, num)
		if err != nil {
			panic("cannot kmins")
		}

		res = make([][]types.Point, len(clusters))
		for i, c := range clusters {
			// solve tsp on all points in cluster
			// check minimum shk bound
			sum_shk = 0
			res[i] = make([]types.Point, len(c))
			for j, idx := range c {
				res[i][j] = pts[idx]
				sum_shk += pts[idx].Shk
			}

			fm := createFastMatrix(res[i])
			_, d := swap2opt(res[i], fm)
			min_shk, _ := fshk(d)

			if min_shk > sum_shk {
				skip = true
				break
			}
		}

		if !skip {
			fmt.Println("clusters: ", len(res))
			return res
		}
		skip = false
	}

	fmt.Println("clusters: ", len(res))
	return res
}

func createFastMatrix(pts []types.Point) [][]float64 {
	res := make([][]float64, len(pts))
	for i := 0; i < len(pts); i++ {
		res[i] = make([]float64, len(pts))
		for j := 0; j < len(pts); j++ {
			if i == j {
				res[i][j] = 0
			} else {
				res[i][j] = CalcSphereDistance(pts[i], pts[j])
			}

		}
	}
	return res
}
