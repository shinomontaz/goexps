package main

import (
	"fmt"
	"time"
)

func main() {
	n := 50
	extraEdges := 4
	maxCap := 5
	seed := time.Now().UnixNano() // для разных запусков будет разный граф

	g := GenerateRandomGraph(n, extraEdges, maxCap, seed)
	g.PrintGraph()

	s, t := PickSourceSink(n, seed+1) // другой seed, чтобы не совпадал с генерацией графа
	fmt.Printf("\nSource: %d, Sink: %d\n", s, t)

	// адгоритм Диница
	maxFlow := Dinic(g, s, t)
	fmt.Printf("Максимальный поток: %d\n", maxFlow)

}
