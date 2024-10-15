package mst

import (
	"lml-drive/types"
	"sort"
)

func forth(pts []types.Point, tm [][]float64, dm [][]float64, dsu *DSU) []int {
	very_big_clusters := map[int]struct{}{}
	max_cluster_size := 150
	max_dur_join := 600.0
	small_cluster_size := 9
	min_cluster_size := 30

	pt_to_cluster := make([]int, len(pts))
	for i := range pts {
		if len(dsu.Cluster(i)) >= max_cluster_size {
			very_big_clusters[dsu.Get(i)] = struct{}{}
		}
		pt_to_cluster[i] = dsu.Get(i)
	}

	edges := []Edge{}
	pt_to_recluster := []types.Point{}
	var unq_points []int
	var best_cluster [2]int
	bad := make([]int, 0)
	sorted_dur := make([][]int, 0)

	maxval := 10000000000000.0

	idxs := make([]int, len(pts))
	for i := range pts {
		idxs[i] = i
	}

	orig_mapping := map[int]int{}
	rev_orig_mapping := map[int]int{}

	for cluster_id := range very_big_clusters {
		unq_points = dsu.Cluster(cluster_id)

		for idx, i := range unq_points {
			orig_mapping[idx] = i
			rev_orig_mapping[i] = idx
			pt_to_recluster = append(pt_to_recluster, pts[i])
			for idx2 := idx + 1; idx2 < len(unq_points); idx2++ {
				j := unq_points[idx2]
				if tm[i+1][j+1] > max_dur_join {
					continue
				}
				edges = append(edges, Edge{from: idx, to: idx2, w: tm[i+1][j+1]})
			}

		}

		bew := ByEdgeWeight{
			list: edges,
		}
		sort.Sort(bew)

		dsu2 := NewDsu(pt_to_recluster, dm)

		for _, e := range edges {
			if dsu2.Get(e.from) == dsu2.Get(e.to) {
				continue
			}

			len1 := len(dsu2.Cluster(e.from))
			len2 := len(dsu2.Cluster(e.to))

			if len1+len2 > max_cluster_size && min(len1, len2) >= small_cluster_size {
				continue
			}

			dsu2.Join(e.from, e.to)
		}

		sorted_dur = sorted_dur[:0]
		for idx, i := range unq_points {
			bclst := ByClosest{list: make([]int, len(unq_points)), src: i, dm: tm}
			copy(bclst.list, unq_points)
			sort.Sort(bclst)
			sorted_dur = append(sorted_dur, make([]int, len(unq_points)))
			copy(sorted_dur[idx], bclst.list)
			for j := 0; j < len(sorted_dur[idx]); j++ {
				sorted_dur[idx][j] = rev_orig_mapping[sorted_dur[idx][j]]
			}
		}

		for {
			bad = bad[:0]

			for i := range unq_points {
				if len(dsu2.Cluster(i)) < min_cluster_size {
					bad = append(bad, dsu2.Get(i))
				}
			}

			if len(bad) == 0 {
				break
			}

			best_cluster = [2]int{-1, -1}
			best_cluster_out_value := maxval

			for _, i := range bad {
				mi := maxval
				mi_id := -1

				for _, j := range dsu2.Cluster(i) {
					for len(sorted_dur[j]) > 0 && dsu2.Get(sorted_dur[j][0]) == dsu2.Get(j) {
						sorted_dur[j] = sorted_dur[j][1:]
					}

					for len(sorted_dur[j]) > 0 {
						if tm[orig_mapping[j]+1][orig_mapping[sorted_dur[j][0]]+1] < mi {
							mi = tm[orig_mapping[j]+1][orig_mapping[sorted_dur[j][0]]+1]
							mi_id = sorted_dur[j][0]
						}
						break
					}

				}

				if best_cluster[0] == -1 || best_cluster_out_value > mi {
					best_cluster[0] = i
					best_cluster[1] = mi_id
					best_cluster_out_value = mi
				}
			}

			if best_cluster[0] == -1 {
				break
			}
			dsu2.Join(best_cluster[0], best_cluster[1])
		}

		for idx, i := range unq_points {
			pt_to_cluster[i] = dsu2.Get(idx)
		}
	}

	return pt_to_cluster
}
