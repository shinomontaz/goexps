package main

import "fmt"

type Bipartite struct {
	left []int
	p    [][]int //
}

func (g *Bipartite) specialDfs(v int, mem map[int]bool, p []int) bool {
	if mem[v] {
		return false
	}
	mem[v] = true

	for _, u := range g.p[v] {
		if p[u] == -1 || g.specialDfs(p[u], mem, p) {
			p[u] = v
			return true
		}
	}
	return false
}

func main() {
	// graph := &Bipartite{
	// 	left: []int{0, 1, 2, 3, 4, 5, 6},
	// 	p:    [][]int{[]int{1}, []int{0}, []int{1, 2, 3}, []int{3}, []int{4}, []int{6}, []int{}},
	// }

	graph := &Bipartite{
		left: []int{0, 1},
		p:    [][]int{[]int{1}, []int{0}},
	}

	n := len(graph.left)
	matching := make([]int, n)
	for idx := range matching {
		matching[idx] = -1
	}

	for i := 0; i < n; i++ {
		mem := make(map[int]bool)
		graph.specialDfs(i, mem, matching)
	}

	for i, m := range matching {
		if m != -1 {
			fmt.Println(m, i)
		}
	}
}
