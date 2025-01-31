package main

type Pair struct {
	idx int
	val float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].val > p[j].val }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
