package bnb

import (
	"fmt"
	"sort"
)

type PriorityQueue struct {
	items [][]interface{}
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{items: [][]interface{}{}}
}

func (pq *PriorityQueue) Push(priority int64, item []int, matrix [][]int64) {
	pq.items = append(pq.items, []interface{}{priority, item, matrix})
	sort.Slice(pq.items, func(i, j int) bool {
		return pq.items[i][0].(int64) < pq.items[j][0].(int64)
	})
}

func (pq *PriorityQueue) Pop() []interface{} {
	item := pq.items[0]
	pq.items = pq.items[1:]
	return item
}

func (pq *PriorityQueue) Size() int {
	return len(pq.items)
}

func (pq *PriorityQueue) Empty() bool {
	return len(pq.items) == 0
}

type TSPstr struct {
	costMatrix [][]int64
	numCities  int
	visited    []bool
	initBound  int64
}

func NewTSPstr(costMatrix [][]int64) *TSPstr {
	numCities := len(costMatrix)
	visited := make([]bool, numCities)
	initBound, costMatrix := reduceMatrix(costMatrix)
	return &TSPstr{costMatrix: costMatrix, numCities: numCities, visited: visited, initBound: initBound}
}

func reduceMatrix(costMatrix [][]int64) (int64, [][]int64) {
	reducedMatrix := make([][]int64, len(costMatrix))
	for i := range costMatrix {
		reducedMatrix[i] = make([]int64, len(costMatrix[i]))
		copy(reducedMatrix[i], costMatrix[i])
	}
	lowBound := int64(0)

	// Subtract the row in cost matrix with the smallest element
	for i := range reducedMatrix {
		minCost := reducedMatrix[i][0]
		for _, cost := range reducedMatrix[i] {
			if cost < minCost {
				minCost = cost
			}
		}
		lowBound += minCost
		for j := range reducedMatrix[i] {
			reducedMatrix[i][j] -= minCost
		}
	}

	// Subtract the col in cost matrix with the smallest element
	for j := 0; j < len(reducedMatrix); j++ {
		minCost := int64(99999)
		for i := range reducedMatrix {
			if reducedMatrix[i][j] < minCost {
				minCost = reducedMatrix[i][j]
			}
		}
		lowBound += minCost
		for i := range reducedMatrix {
			reducedMatrix[i][j] -= minCost
		}
	}

	return lowBound, reducedMatrix
}

func (tsp *TSPstr) toInfinite(row, col int, matrix [][]int64) [][]int64 {
	for j := range matrix[row] {
		matrix[row][j] = 99999
	}
	for i := range matrix {
		matrix[i][col] = 99999
	}
	matrix[col][row] = 99999
	return matrix
}

func (tsp *TSPstr) lowerBound(lastTour, currentTour []int, lastBound int64) (int64, [][]int64) {
	costMatrix := make([][]int64, len(tsp.costMatrix))
	for i := range tsp.costMatrix {
		costMatrix[i] = make([]int64, len(tsp.costMatrix[i]))
		copy(costMatrix[i], tsp.costMatrix[i])
	}

	nextCity := currentTour[len(currentTour)-1]
	currentCity := lastTour[len(lastTour)-1]

	cost := costMatrix[currentCity][nextCity]

	for j := range costMatrix[currentCity] {
		costMatrix[currentCity][j] = 99999
	}

	for i := range costMatrix {
		costMatrix[i][nextCity] = 99999
	}

	costMatrix[nextCity][currentCity] = 99999

	_, costMatrix = reduceMatrix(costMatrix)
	costMatrix = tsp.toInfinite(currentCity, nextCity, costMatrix)

	return cost, costMatrix
}

func (tsp *TSPstr) branchAndBound() ([]int, int64, [][]int64) {
	prioQueue := NewPriorityQueue()
	initialTour := []int{0}
	tsp.visited[0] = true
	initialBound := tsp.initBound
	prioQueue.Push(initialBound, initialTour, tsp.costMatrix)

	bestTour := []int{}
	bestCost := int64(99999)
	bestCostMatrix := [][]int64{}

	for !prioQueue.Empty() {
		currentNode := prioQueue.Pop()
		currentBound := currentNode[0].(int64)
		currentTour := currentNode[1].([]int)
		currentCostMatrix := currentNode[2].([][]int64)

		if currentBound >= bestCost {
			continue
		}

		if len(currentTour) == tsp.numCities {
			if currentBound < bestCost {
				bestTour = currentTour
				bestCost = currentBound
				bestCostMatrix = currentCostMatrix
			}
			continue
		}

		for _, city := range currentTour {
			tsp.visited[city] = true
		}

		for j := 0; j < tsp.numCities; j++ {
			if !tsp.visited[j] {
				newTour := append(currentTour, j)
				newBound, newCostMatrix := tsp.lowerBound(currentTour, newTour, currentBound)
				prioQueue.Push(newBound+currentBound, newTour, newCostMatrix)
			}
		}

		for _, city := range currentTour {
			tsp.visited[city] = false
		}
	}

	return bestTour, bestCost, bestCostMatrix
}

func printMatrix(matrix [][]int64) {
	for _, row := range matrix {
		fmt.Println(row)
	}
}
