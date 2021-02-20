package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGraph()
	n := 10
	for i := 1; i < n; i++ {
		g.AddNode(i)
	}
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			dice := rand.Float64()
			if dice > 0.3 {
				dist := 10.0 * float64(rand.Intn(10)) * rand.Float64()
				g.AddEdge(i, j, dist)
			}
		}
	}

	start := rand.Intn(n)

	paths := g.Dijkstra(start)

	fmt.Println("start node:", start)
	for _, path := range paths {
		path.Print()
	}
}
