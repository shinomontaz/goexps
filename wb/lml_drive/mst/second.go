package mst

import (
	"lml-drive/types"
	"math"
	"sort"
)

func second(pts []types.Point, tm [][]float64, dsu *DSU, flag bool) {
	max_cluster_size := 150

	sorted_dist := make([][]int, len(pts))
	idxs := make([]int, len(pts))
	for i := range pts {
		idxs[i] = i
	}

	for i := range pts {
		bclst := ByClosest{list: make([]int, len(pts)), src: i, dm: tm}
		copy(bclst.list, idxs)
		sort.Sort(bclst)
		sorted_dist[i] = make([]int, len(pts))
		copy(sorted_dist[i], bclst.list)
	}

	bad := make([]int, 0, len(pts))
	for {
		bad = bad[:0]
		for i := range pts {
			if dsu.GetShk(i) >= minShk(dsu.GetUpper(i)-1) || len(sorted_dist[i]) == 0 {
				continue
			}
			bad = append(bad, i)
		}

		best_join := [2]int{-1, -1}
		best_join_dist := 100000000000000000000.0

		for _, i := range bad {
			for len(sorted_dist[i]) > 0 && dsu.Get(sorted_dist[i][0]) == dsu.Get(i) {
				sorted_dist[i] = sorted_dist[i][1:]
			}

			if len(sorted_dist[i]) == 0 {
				continue
			}

			pt_dist := tm[i+1][sorted_dist[i][0]+1]
			if pt_dist >= best_join_dist {
				continue
			}

			if flag && len(dsu.Cluster(sorted_dist[i][0])) >= max_cluster_size {
				continue
			}

			bucket1, bucket2 := dsu.GetIdx(i), dsu.GetIdx(sorted_dist[i][0])

			if math.Abs(float64(bucket1)-float64(bucket2)) > 1 {
				sorted_dist[i] = sorted_dist[i][:0]
				continue
			}
			if bucket1 < bucket2 {
				continue
			}

			if pt_dist < best_join_dist {
				best_join[0] = i
				best_join[1] = sorted_dist[i][0]
				best_join_dist = pt_dist
			}
		}

		if best_join[0] == -1 {
			break
		}
		dsu.Join(best_join[0], best_join[1])
	}
}
