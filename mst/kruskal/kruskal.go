package main

import "sort"

type edge struct {
	U, V   int
	Weight float64
}

func KruskalMst(n int, edges []edge) ([]int, float64) {
	// Initialize variables to store the minimum spanning tree and its total cost
	mst := []int{}
	var cost float64

	// Create a new UnionFind data structure with 'n' nodes
	u := NewUnionFind(n)

	// Sort the edges in non-decreasing order based on their weights
	sort.SliceStable(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	// Iterate through the sorted edges
	for _, edge := range edges {
		// Check if adding the current edge forms a cycle or not
		if u.Find(edge.U) != u.Find(edge.V) {
			// Add the edge to the minimum spanning tree
			//			mst = append(mst, edge)
			mst = append(mst, edge.U, edge.V)
			// Add the weight of the edge to the total cost
			cost += edge.Weight
			// Merge the sets containing the start and end vertices of the current edge
			u.Union(edge.U, edge.V)
		}
	}

	// Return the minimum spanning tree and its total cost
	return mst, cost
}

func MyMst(n int, edges []edge) ([]int, float64) {
	// Sort edges by weight
	sort.SliceStable(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	// Union-Find data structures
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
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
			} else if rank[rootX] < rank[rootY] {
				parent[rootX] = rootY
			} else {
				parent[rootY] = rootX
				rank[rootX]++
			}
		}
	}

	mst := []int{}
	cost := 0.0
	for _, edge := range edges {
		if find(edge.U) != find(edge.V) {
			union(edge.U, edge.V)
			mst = append(mst, edge.U, edge.V)
			cost += edge.Weight
		}
	}

	return mst, cost
}

func buildMst(pts []Point, dm [][]float64) ([]int, float64) {
	n := len(pts)
	edges := []edge{}

	// Create a list of all edges
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, edge{i, j, dm[i][j]})
		}
	}

	return MyMst(n, edges)
}
