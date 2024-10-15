package mst

import "lml-drive/types"

type DSU struct {
	p       []int // для каждой точки - номер ее кластера
	d       []float64
	shk_sum []int         // для каждого кластера - сумма ШК в нем
	dsu     map[int][]int // для каждого кластера - список точек в нем
	bucket  []int         // индекс ограничения
	upper   []float64     // верхняя граница ограничения
}

func NewDsu(pts []types.Point, dm [][]float64) *DSU {
	dsu := DSU{
		p:       make([]int, len(pts)),
		d:       make([]float64, len(pts)),
		shk_sum: make([]int, len(pts)),
		dsu:     make(map[int][]int, len(pts)),
		bucket:  make([]int, len(pts)),
		upper:   make([]float64, len(pts)),
	}

	for i, pt := range pts {
		dsu.p[i] = i
		dsu.d[i] = dm[0][i+1] + dm[i+1][0]
		dsu.shk_sum[i] = pt.Shk
		dsu.dsu[i] = []int{i}
		dsu.bucket[i] = getBuc((dm[0][i+1] + dm[i+1][0]) * 1000.0)
		dsu.upper[i] = float64(maxShk(dm[0][i+1] + dm[i+1][0]*1000)) // p.GetUpperDist(dm[0][i+1] + dm[i+1][0] * 1000)
	}

	return &dsu
}

func (dsu *DSU) Join(a, b int) {
	a = dsu.p[a]
	b = dsu.p[b]

	if a == b {
		return
	}

	if len(dsu.dsu[a]) > len(dsu.dsu[b]) {
		a, b = b, a
	}

	dsu.p[a] = b
	dsu.shk_sum[b] += dsu.shk_sum[a]

	if dsu.d[a] > dsu.d[b] {
		dsu.d[b] = dsu.d[a]
	}

	if dsu.upper[a] > dsu.upper[b] {
		dsu.bucket[b] = dsu.bucket[a]
		dsu.upper[b] = dsu.upper[a]
	}

	for _, i := range dsu.dsu[a] {
		dsu.p[i] = b
	}

	dsu.dsu[b] = append(dsu.dsu[b], dsu.dsu[a]...)
	delete(dsu.dsu, a)
}

func (dsu *DSU) Get(a int) int {
	return dsu.p[a]
}

func (dsu *DSU) GetIdx(i int) int {
	return dsu.bucket[dsu.p[i]]
}

func (dsu *DSU) GetShk(i int) int {
	return dsu.shk_sum[dsu.p[i]]
}

func (dsu *DSU) GetUpper(i int) float64 {
	return dsu.upper[dsu.p[i]]
}

func (dsu *DSU) Cluster(i int) []int {
	return dsu.dsu[dsu.p[i]]
}
