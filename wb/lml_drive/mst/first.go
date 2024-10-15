package mst

import "lml-drive/types"

type item struct {
	i int
	w float64
}

func first(pts []types.Point, tm [][]float64, dsu *DSU) {
	used := map[int]bool{}
	max_dur_join := 420.0
	min_cluster_size := 30

	max_v := 1000000000000.0
	best_edge_w := max_v

	mi_in := make([]float64, len(pts))
	in_id := make([]int, len(pts))

	qlist := []item{}

	for len(used) < len(pts) {
		best_edge_w = max_v
		best_v := -1
		for i := range pts {
			mi_in[i] = max_v
			in_id[i] = 0
			if used[i] {
				continue
			}
			for j := range pts {
				if i == j {
					continue
				}
				if best_v == -1 || tm[i+1][j+1] < best_edge_w {
					best_v = i // здесь не надо i+1, но +1 надо будет учесть потом
					best_edge_w = tm[i+1][j+1]
				}
			}
		}
		if best_v == -1 || best_edge_w > max_dur_join {
			break
		}

		qlist = qlist[:0]
		qlist = append(qlist, item{i: best_v, w: best_edge_w})
		in_id[best_v] = best_v

		for len(qlist) > 0 {
			v := GetMin(qlist)
			qlist = Discard(qlist, v)
			used[v] = true
			dsu.Join(v, in_id[v]) // в DSU нумерация с 0, где 0 это первая точка, а не склад
			for j := range pts {
				if len(dsu.Cluster(v)) >= min_cluster_size && used[j] { //
					continue
				}
				if tm[v+1][j+1] < mi_in[j] {
					qlist = Discard(qlist, j)
					mi_in[j] = tm[v+1][j+1]
					in_id[j] = v
					if tm[v+1][j+1] <= max_dur_join {
						qlist = append(qlist, item{i: j, w: mi_in[j]})
					}
				}
			}

		}

	}

}

func GetMin(qlist []item) int {
	min := qlist[0].w
	minIdx := qlist[0].i

	for _, it := range qlist {
		if it.w < min {
			min = it.w
			minIdx = it.i
		}

	}

	return minIdx
}

func Discard(qlist []item, i int) []item {
	for idx, it := range qlist {
		if it.i == i {
			return append(qlist[:idx], qlist[idx+1:]...)
		}
	}

	return qlist
}
