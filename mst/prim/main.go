package main

import (
	"fmt"
	"log"
	"math"
)

var dm = [][]float64{
	{0, 44203, 43854, 42716, 40774, 44084, 42973, 45933, 41916},
	{44374, 0, 349, 1486, 5280, 1206, 1597, 1557, 2036},
	{44025, 349, 0, 1137, 4931, 857, 1247, 1965, 1687},
	{43162, 1486, 1137, 0, 4016, 1368, 256, 3216, 772},
	{41406, 5280, 4931, 4016, 0, 5161, 4127, 7010, 3453},
	{44255, 1206, 857, 1368, 5161, 0, 1478, 2099, 1918},
	{43220, 1597, 1247, 256, 4127, 1478, 0, 3326, 883},
	{46104, 1557, 1965, 3216, 7010, 2099, 3326, 0, 3766},
	{42547, 2036, 1687, 772, 3453, 1918, 883, 3766, 0},
}

func main() {
	// берем произвольную точку
	// из этой точки беремминимальное ребро
	// добавляем его в остов
	// из новых точек берем минимальное ребро
	// повторяем пока не покроем все точки

	mst := make([]int, 0, len(dm))
	mst_edges := make([][2]int, 0, len(dm))
	visited := make(map[int]bool)

	start := 0
	mst = append(mst, start)
	visited[start] = true
	for len(mst) < len(dm) {
		minDist := float64(math.MaxInt32)
		var from, to int
		var foundMin bool
		for _, v := range mst {
			for i := 0; i < len(dm); i++ {
				if visited[i] || i == v {
					continue
				}
				if dm[v][i] < minDist {
					minDist = dm[v][i]
					from = v
					to = i
					foundMin = true
				}
			}
		}

		if !foundMin {
			log.Fatal("disconnected graph")
		}

		mst = append(mst, to)
		mst_edges = append(mst_edges, [2]int{from, to})
		visited[to] = true
	}

	fmt.Println(mst)
}
