package main

type UnionFind struct {
	parent []int
	rank   []int
}

// Initialise a new union find data structure with s nodes
func NewUnionFind(s int) UnionFind {
	parent := make([]int, s)
	rank := make([]int, s)
	for i := 0; i < s; i++ {
		parent[i] = i
		rank[i] = 1
	}
	return UnionFind{parent, rank}
}

// Find finds the root of the set to which the given element belongs.
// It performs path compression to make future Find operations faster.
func (u *UnionFind) Find(q int) int {
	if q != u.parent[q] {
		u.parent[q] = u.Find(u.parent[q])
	}
	return u.parent[q]
}

// Union merges the sets, if not already merged, to which the given elements belong.
// It performs union by rank to keep the tree as flat as possible.
func (u *UnionFind) Union(p, q int) {
	rootP := u.Find(p)
	rootQ := u.Find(q)

	if rootP == rootQ {
		return
	}

	if u.rank[rootP] < u.rank[rootQ] {
		u.parent[rootP] = rootQ
	} else if u.rank[rootP] > u.rank[rootQ] {
		u.parent[rootQ] = rootP
	} else {
		u.parent[rootQ] = rootP
		u.rank[rootP]++
	}
}
