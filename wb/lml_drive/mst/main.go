package mst

import "lml-drive/types"

type int_from_dist_f func(d float64) int

var minShk int_from_dist_f
var maxShk int_from_dist_f
var getBuc int_from_dist_f

func Init(mnshk, mxshk, bushk int_from_dist_f) {
	minShk = mnshk
	maxShk = mxshk
	getBuc = bushk
}

func Do(pts []types.Point, dm, tm [][]float64) [][]types.Point {
	res := [][]types.Point{}

	dsu := NewDsu(pts[1:], dm)
	first(pts[1:], tm, dsu)
	second(pts[1:], tm, dsu, true)
	second(pts[1:], tm, dsu, false)
	third(pts[1:], tm, dsu)
	pt_to_cluster := forth(pts[1:], tm, dm, dsu)

	clusters := map[int]int{}
	for p, c := range pt_to_cluster {
		if _, ok := clusters[c]; !ok {
			clusters[c] = len(res)
			res = append(res, []types.Point{pts[0]}) // первой точкой в наборе всегда будет склад
		}

		cl := clusters[c]

		res[cl] = append(res[cl], pts[p+1]) // потому что на кластеризацию подавали pts[1:] и индексы в ответе начинаются с 0
	}
	return res
}
