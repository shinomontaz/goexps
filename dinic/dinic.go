package main

import "math"

func Dinic(g *Graph, s, t int) (maxFlow int) {
	// Implementation of Dinic's algorithm

	flows := make([][]int, len(g.Edges))
	for i := range flows {
		flows[i] = make([]int, len(g.Edges))
	}

	for {
		levels := bfs(g, s, flows)
		if levels[t] == -1 {
			return maxFlow
		}

		for {
			mem := make(map[int]bool)
			mf := dfs(g, s, t, math.MaxInt, levels, flows, mem)
			if mf == 0 {
				break
			}

			maxFlow += mf
		}

	}
}

func bfs(g *Graph, s int, flows [][]int) []int {
	// запускаем BFS, строим уровни
	seen := map[int]struct{}{}
	queue := []int{s}

	levels := make([]int, len(g.Edges))
	for i := range len(levels) {
		levels[i] = -1
	}

	levels[s] = 0
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		seen[u] = struct{}{}

		for _, e := range g.Edges[u] {
			if _, ok := seen[e.To]; !ok {
				if e.Cap-flows[u][e.To] > 0 {
					levels[e.To] = levels[u] + 1
					queue = append(queue, e.To)
					seen[e.To] = struct{}{}
				}
			}
		}
	}

	return levels
}

func dfs(g *Graph, n, t, flow int, levels []int, flows [][]int, mem map[int]bool) int {
	// мы може
	if n == t || flow == 0 {
		return flow
	}

	mem[n] = true

	for _, v := range g.Edges[n] {
		if levels[v.To] == levels[n]+1 && !mem[v.To] { // только на следующий уровень
			currFlow := dfs(g, v.To, t, min(flow, v.Cap-flows[n][v.To]), levels, flows, mem)
			if currFlow > 0 {
				flows[n][v.To] += flow
				flows[v.To][n] -= flow

				return currFlow
			}
		}
	}

	return 0
}
