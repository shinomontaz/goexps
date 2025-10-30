package bnb

import (
	"fmt"
	"math"
)

func firstMin(adj [][]int64, i int) int64 {
	min := int64(math.MaxInt64)
	for k := 0; k < len(adj); k++ {
		if adj[i][k] < min && i != k {
			min = adj[i][k]
		}
	}
	return min
}

func secondMin(adj [][]int64, i int) int64 {
	first, second := int64(math.MaxInt64), int64(math.MaxInt64)
	for j := 0; j < len(adj); j++ {
		if i == j {
			continue
		}
		if adj[i][j] <= first {
			second = first
			first = adj[i][j]
		} else if adj[i][j] <= second && adj[i][j] != first {
			second = adj[i][j]
		}
	}
	return second
}

func TSPRec(adj [][]int64, currBound int64, currWeight int64, level int, currPath []int, start int, last int, visited []bool, finalPath []int, finalRes *int64) {
	if level == len(adj) {
		if last == -1 {
			currRes := currWeight
			if currRes < *finalRes {
				copy(finalPath, currPath)
				*finalRes = currRes
			}
		} else if currPath[level-1] == last {
			currRes := currWeight
			if currRes < *finalRes {
				copy(finalPath, currPath)
				*finalRes = currRes
			}
		}
		return
	}

	for i := 0; i < len(adj); i++ {
		if adj[currPath[level-1]][i] != 0 && !visited[i] {
			temp := currBound
			currWeight += adj[currPath[level-1]][i]

			if level == 1 {
				currBound -= (firstMin(adj, currPath[level-1]) + firstMin(adj, i)) / 2
			} else {
				currBound -= (secondMin(adj, currPath[level-1]) + firstMin(adj, i)) / 2
			}

			if currBound+currWeight < *finalRes {
				currPath[level] = i
				visited[i] = true
				TSPRec(adj, currBound, currWeight, level+1, currPath, start, last, visited, finalPath, finalRes)
			}

			currWeight -= adj[currPath[level-1]][i]
			currBound = temp

			for j := 0; j < len(adj); j++ {
				visited[j] = false
			}
			for j := 0; j <= level-1; j++ {
				visited[currPath[j]] = true
			}
		}
	}
}

func TSP(adj [][]int64, start int, last int) ([]int, int64) {
	n := len(adj)
	currPath := make([]int, n)
	for i := range currPath {
		currPath[i] = -1
	}

	currBound := int64(0)
	for i := 0; i < n; i++ {
		currBound += firstMin(adj, i) + secondMin(adj, i)
	}

	currBound = (currBound & 1) / 2
	if currBound&1 != 0 {
		currBound++
	}

	visited := make([]bool, n)
	finalPath := make([]int, n)
	finalRes := int64(math.MaxInt64)

	visited[start] = true
	currPath[0] = start
	TSPRec(adj, currBound, 0, 1, currPath, start, last, visited, finalPath, &finalRes)

	return finalPath, finalRes
}

func main() {
	adj := [][]int64{
		{math.MaxInt64, 10, 15, 20, 25},
		{10, math.MaxInt64, 35, 25, 30},
		{15, 35, math.MaxInt64, 30, 20},
		{20, 25, 30, math.MaxInt64, 15},
		{25, 30, 20, 15, math.MaxInt64},
	}
	startPoint := 0
	lastPoint := -1

	path, cost := TSP(adj, startPoint, lastPoint)
	fmt.Printf("Minimum cost: %d\n", cost)
	fmt.Printf("Path Taken: ")
	for i := 0; i < len(path); i++ {
		fmt.Printf("%d ", path[i])
	}
	fmt.Println()
}
