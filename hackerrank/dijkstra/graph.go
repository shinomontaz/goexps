package main

import (
	"fmt"
	"math"
	"sort"
)

type Node struct {
	id int
}

type Edge struct {
	start *Node
	end   *Node
	cost  float64
}

type Graph struct {
	edges []*Edge
	nodes map[int]*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int]*Node),
		edges: make([]*Edge, 0),
	}
}

func (g *Graph) AddNode(idx int) {
	g.nodes[idx] = &Node{id: idx}
}

func (g *Graph) AddEdge(idx1, idx2 int, cost float64) {
	// check id existance
	if _, exists := g.nodes[idx1]; !exists {
		g.AddNode(idx1)
	}
	if _, exists := g.nodes[idx2]; !exists {
		g.AddNode(idx2)
	}

	n1 := g.nodes[idx1]
	n2 := g.nodes[idx2]

	g.edges = append(g.edges, &Edge{
		start: n1,
		end:   n2,
		cost:  cost,
	})
}

func getClosestNonvisited(distances map[*Node]float64, visited []*Node) *Node {
	min := 0.0
	var closest *Node

	for n, dist := range distances {
		isVisited := false
		for _, v := range visited {
			if v == n {
				isVisited = true
				break
			}
		}

		if !isVisited {
			if min < dist {
				closest = n
				min = dist
			}
		}
	}

	return closest
}

func (g *Graph) getEdges(n *Node) []*Edge {
	edges := make([]*Edge, 0)

	for _, e := range g.edges {
		if e.start == n {
			edges = append(edges, e)
		}
	}

	return edges
}

func (g *Graph) Print() {
	matrix := make([][]float64, len(g.nodes))

	ids := make([]int, 0, len(g.nodes))
	for id := range g.nodes {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for ii, i := range ids {
		n := g.nodes[i]
		matrix[ii] = make([]float64, len(g.nodes))

		for jj, j := range ids {
			m := g.nodes[j]
			if ii == jj {
				matrix[ii][jj] = 0
			}

			for _, e := range g.edges {
				if e.start == n && e.end == m {
					matrix[ii][jj] = e.cost
				}
			}
		}
	}

	for _, row := range matrix {
		rowstr := make([]string, 0, len(g.nodes))
		for _, cost := range row {
			rowstr = append(rowstr, fmt.Sprintf("%f", cost))
		}
		fmt.Println(rowstr)
	}
}

type Path struct {
	dest  *Node
	cost  float64
	steps []*Node
}

func (p *Path) Print() {
	fmt.Printf("Dest: %d, Cost: %f\n", p.dest.id, p.cost)
}

func (g *Graph) Dijkstra(startId int) map[*Node]Path {
	var start *Node
	for _, n := range g.nodes {
		if n.id == startId {
			start = n
			break
		}
	}

	res := make(map[*Node]Path)

	distances := make(map[*Node]float64)
	distances[start] = 0
	for _, n := range g.nodes {
		distances[n] = math.Inf(1)
	}

	visited := make([]*Node, 0)

	for len(visited) < len(g.nodes) {
		n := getClosestNonvisited(distances, visited) // this node connected to start node

		visited = append(visited, n)

		edges := g.getEdges(n) // this edges have start = n node from getClosestNonvisited

		if math.IsInf(distances[n], 1) {
			distances[n] = 0
		}

		for _, e := range edges {
			dist := distances[n] + e.cost
			if dist < distances[e.end] {
				distances[e.end] = dist // релаксация
			}
		}
	}

	for n, cost := range distances {
		res[n] = Path{
			dest: n,
			cost: cost,
		}
	}

	return res
}
