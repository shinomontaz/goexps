package main

type Path struct {
	way        []int
	points     []*LatLng
	neighbours map[int][]int
}

func (p *Path) Index(i int) int {
	//	Return the index of a node in a tour.

	for idx, it := range p.way {
		if it == i {
			return idx
		}
	}

	panic("not existed node index requested")
}

func (p *Path) Around(node int) []int {
	/*
	   Return the predecessor and successor of the current node, given by
	   index.

	   Parameters:

	       - node: node to look around

	   Return: (pred, succ)
	*/
	index := p.index(node)

	pred := index - 1
	succ := index + 1

	if succ == len(p.way) {
		succ = 0
	}

	return []int{p.way[pred], p.way[succ]}
}

func (p *Path) pred(index int, prev bool) {

	//	Return the predecessor or successor depending on the `prev` parameter.

	if prev {
		return p.way[index-1]
	}

	return p.way[index+1]

}
