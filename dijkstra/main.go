package main

import (
	"fmt"
)

// a time has come: implement dijkstra myself
func main() {
	a := &Node{Id: 1}
	b := &Node{Id: 2}
	c := &Node{Id: 3}
	d := &Node{Id: 4}
	e := &Node{Id: 5}
	f := &Node{Id: 6}
	g := &Node{Id: 7}

	graph := &Graph{}
	graph.AddEdge(a, c, 2)
	graph.AddEdge(a, b, 5)
	graph.AddEdge(c, b, 1)
	graph.AddEdge(c, d, 9)
	graph.AddEdge(b, d, 4)
	graph.AddEdge(d, e, 2)
	graph.AddEdge(d, g, 30)
	graph.AddEdge(d, f, 10)
	graph.AddEdge(f, g, 1)

	res := dijkstra(graph, a)

	for _, r := range res {
		fmt.Println(r.toString())
	}
}
