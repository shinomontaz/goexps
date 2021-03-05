package main

import (
	"fmt"
	"testing"
)

func TestAlgo(t *testing.T) {
	graph := NewGraph()
	graph.AddEdge(1, 3, 2)
	graph.AddEdge(1, 2, 5)
	graph.AddEdge(3, 2, 1)
	graph.AddEdge(3, 4, 9)
	graph.AddEdge(2, 4, 4)
	graph.AddEdge(4, 5, 2)
	graph.AddEdge(4, 7, 30)
	graph.AddEdge(4, 6, 10)
	graph.AddEdge(6, 7, 1)

	paths := graph.Dijkstra(1)

	graph.Print()
	fmt.Println("start node:", 1)
	for _, path := range paths {
		path.Print()
	}
}
