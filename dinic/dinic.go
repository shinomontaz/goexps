package main

import "fmt"

func Dinic(g *Graph, s, t int) (maxFlow int) {
	// Implementation of Dinic's algorithm

	// запускаем BFS, строим уровни
	seen := map[int]struct{}{}
	queue := []int{s}

	levels := make([]int, len(g.Edges))
	for i := 0; i < len(levels); i++ {
		levels[i] = -1
	}

	levels[s] = 0
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		seen[u] = struct{}{}

		for _, e := range g.Edges[u] {
			if _, ok := seen[e.To]; !ok {
				if e.Cap > 0 {
					levels[e.To] = levels[u] + 1
					queue = append(queue, e.To)
					seen[e.To] = struct{}{}
				}
			}
		}
	}

	fmt.Println(levels)

	return 0
}
