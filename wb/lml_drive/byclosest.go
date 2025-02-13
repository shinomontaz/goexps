package main

import "lml-drive/types"

type ByClosest struct {
	list []int
	pts  []types.Point
	tm   [][]float64
}

func (a ByClosest) Len() int      { return len(a.list) }
func (a ByClosest) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a ByClosest) Less(i, j int) bool {
	return a.tm[0][a.list[i]] < a.tm[0][a.list[j]]
}
