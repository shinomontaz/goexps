package main

type ByClosest struct {
	pts []Point
	tm  [][]float64
}

func (a ByClosest) Len() int      { return len(a.pts) }
func (a ByClosest) Swap(i, j int) { a.pts[i], a.pts[j] = a.pts[j], a.pts[i] }
func (a ByClosest) Less(i, j int) bool {
	return a.tm[0][a.pts[i].Cid] < a.tm[0][a.pts[j].Cid]
}
