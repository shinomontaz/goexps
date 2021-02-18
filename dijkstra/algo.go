package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Path struct {
	Dest *Node    // id of a target node
	Cost float64  // total travel cost
	Path []string // orderd slice of a nodes to reach destination node
}

func (p *Path) toString() string {
	return fmt.Sprintf("Distance to %d = %f (path: %s))\n", p.Dest.Id, p.Cost, strings.Join(p.Path, ",  "))
}

func (g *Graph) NewCostMatrix(start *Node) map[*Node]float64 {
	res := make(map[*Node]float64)

	res[start] = 0

	for _, n := range g.nodes {
		if n != start {
			res[n] = math.Inf(1)
		}
	}

	return res
}

func NodeInList(n *Node, list []*Node) bool {
	for _, v := range list {
		if v == n {
			return true
		}
	}

	return false
}

func getClosestNonvisited(m map[*Node]float64, visited []*Node) *Node {
	paths := []Path{}
	for n, cost := range m {
		if !NodeInList(n, visited) {
			paths = append(paths, Path{
				Dest: n,
				Cost: cost,
			})
		}
	}

	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Cost < paths[j].Cost
	})

	return paths[0].Dest
}

func (g *Graph) NodeEdges(n *Node) []*Edge {
	edges := []*Edge{}
	for _, e := range g.edges {
		if e.Start == n {
			edges = append(edges, e)
		}
	}
	return edges
}

func dijkstra(g *Graph, start *Node) []Path {
	res := []Path{}

	ancestors := make(map[*Node]*Node)

	matrix := g.NewCostMatrix(start)
	visited := []*Node{}

	for len(visited) < len(g.nodes) {
		// loop
		candidate := getClosestNonvisited(matrix, visited)
		visited = append(visited, candidate)

		edges := g.NodeEdges(candidate)

		for _, e := range edges {
			if math.IsInf(matrix[candidate], 1) {
				matrix[candidate] = 0
			}
			dist := matrix[candidate] + e.Cost
			if dist < matrix[e.End] {
				matrix[e.End] = dist
				ancestors[e.End] = candidate
			}
		}
	}

	for n, cost := range matrix {
		path := []string{}
		v := n
		for v != start {
			ancestor := ancestors[v]
			path = append(path, fmt.Sprintf("%d -> %d", ancestor.Id, v.Id))
			v = ancestor
		}

		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}

		res = append(res, Path{
			Dest: n,
			Cost: cost,
			Path: path,
		})
	}

	return res
}
