package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Id int
}

type Edge struct {
	Start *Node
	End   *Node
	Cost  float64
}

type Graph struct {
	nodes []*Node
	edges []*Edge
}

func (g *Graph) AddNode(n *Node) {
	// check existence
	for _, nd := range g.nodes {
		if nd.Id == n.Id {
			return
		}
	}

	g.nodes = append(g.nodes, n)
}

func (g *Graph) AddEdge(start, end *Node, cost float64) {
	g.edges = append(g.edges, &Edge{
		Start: start,
		End:   end,
		Cost:  cost,
	})

	g.AddNode(start)
	g.AddNode(end)
}

func (g *Graph) toString() string {
	s := ""

	s += "Edges:\n"
	for _, e := range g.edges {
		s += fmt.Sprintf("%d -> %d = %f\n", e.Start.Id, e.End.Id, e.Cost)
	}
	s += "\n"
	tmpStringSlice := []string{}
	for _, e := range g.nodes {
		tmpStringSlice = append(tmpStringSlice, fmt.Sprintf("%d", e.Id))
	}

	s += fmt.Sprintf("%s\n", strings.Join(tmpStringSlice, ","))
	return s
}
