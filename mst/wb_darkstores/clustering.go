package main

import "sort"

type edge struct {
	U, V   int
	Weight float64
}

func mstClusterize(clusters, n int, edges []edge, cluster_num [2]int) ([][]int, float64) {
	// Sort edges by weight
	sort.SliceStable(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	penalty = edges[len(edges)-1].Weight

	// Union-Find data structures
	parent := make([]int, n)
	rank := make([]int, n)
	numbers := make(map[int][]int) // для каждого кластера все точки в нем
	for i := range parent {
		parent[i] = i
		numbers[i] = []int{i} // каждый кластер изначально состоит из одной точки
	}

	// Find function with path compression
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	// Union function with union by rank
	union := func(x, y int) {
		rootX := find(x)
		rootY := find(y)

		if rootX != rootY {
			if rank[rootX] > rank[rootY] {
				parent[rootY] = rootX
				numbers[rootX] = append(numbers[rootX], numbers[rootY]...) // объединяем кластеры
				delete(numbers, rootY)                                     // очищаем старый кластер
			} else if rank[rootX] < rank[rootY] {
				parent[rootX] = rootY
				numbers[rootY] = append(numbers[rootY], numbers[rootX]...) // объединяем кластеры
				delete(numbers, rootX)
			} else {
				parent[rootY] = rootX
				rank[rootX]++
				numbers[rootX] = append(numbers[rootX], numbers[rootY]...) // объединяем кластеры
				delete(numbers, rootY)                                     // очищаем старый кластер
			}
		}
	}

	cost := 0.0
	for _, edge := range edges {
		rootU := find(edge.U)
		rootV := find(edge.V)
		if rootU != rootV && len(numbers[rootU])+len(numbers[rootV]) <= cluster_num[1] {
			union(edge.U, edge.V)
			cost += edge.Weight
		}
	}

	res := make([][]int, len(numbers))
	i := 0
	for _, cluster := range numbers {
		res[i] = make([]int, len(cluster))
		copy(res[i], cluster)
		i++
	}

	return res, cost
}

func buildMst(pts []Point, dm [][]float64) ([][]int, float64) {
	n := len(pts)
	edges := []edge{}

	// Create a list of all edges
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, edge{i, j, dm[i+1][j+1]})
		}
	}

	return mstClusterize(clusters, n, edges, cluster_num)
}
