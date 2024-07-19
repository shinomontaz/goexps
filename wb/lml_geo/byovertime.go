package main

type ByOvertime struct {
	list   []int
	routes [][]int
	tm     [][]float64
}

func (a ByOvertime) Len() int      { return len(a.list) }
func (a ByOvertime) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a ByOvertime) Less(i, j int) bool {
	_, oi := cost(a.routes[a.list[i]], a.tm)
	_, oj := cost(a.routes[a.list[j]], a.tm)
	return oi > oj
}
