package mst

type ByClosest struct {
	list []int
	src  int
	dm   [][]float64
}

func (a ByClosest) Len() int      { return len(a.list) }
func (a ByClosest) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a ByClosest) Less(i, j int) bool {
	return a.dm[a.src+1][a.list[i]+1] < a.dm[a.src+1][a.list[j]+1]
}

type ByEdgeWeight struct {
	list []Edge
}

func (a ByEdgeWeight) Len() int      { return len(a.list) }
func (a ByEdgeWeight) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a ByEdgeWeight) Less(i, j int) bool {
	return a.list[i].w < a.list[j].w
}
