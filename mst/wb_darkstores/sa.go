package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var rng *rand.Rand
var penalty float64

func sa(initSolution [][]int, dm [][]float64) ([][]int, float64) {

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	T := 1.0
	cooling := 0.9999
	Tmin := 0.01

	var newSolution, currSolution [][]int
	currSolution = make([][]int, len(initSolution))
	for i, r := range initSolution {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}
	newSolution = make([][]int, len(currSolution))
	oldFitness := fitness(currSolution, dm)
	fitTreshold := oldFitness

	fmt.Println(oldFitness)

	for T > Tmin {
		newSolution = newSolution[:len(currSolution)]
		copySlice(newSolution, currSolution)

		newSolution = mutate(newSolution)

		if !isFeasible(newSolution) {
			continue
		}

		newFitness := fitness(newSolution, dm)
		if newFitness < oldFitness {
			currSolution = currSolution[:len(newSolution)]
			copySlice(currSolution, newSolution)
			oldFitness = newFitness
		} else if newFitness <= fitTreshold && math.Exp((oldFitness-newFitness)/T) > rng.Float64() {
			currSolution = currSolution[:len(newSolution)]
			copySlice(currSolution, newSolution)
			oldFitness = newFitness
		}

		T *= cooling
	}

	fmt.Println(oldFitness)

	return currSolution, oldFitness

}

func mutate(mutated [][]int) [][]int {
	dice := rng.Float64()
	if dice <= 0.25 {
		return swapmove(mutated)
	} else if dice > 0.25 && dice <= 0.75 {
		return insertmove(mutated)
	}

	return twooptmove(mutated)
}

func fitness(in [][]int, dm [][]float64) float64 {
	c := 0.0
	for _, r := range in {
		c += cost(r, dm)
	}
	return c
}

func cost(route []int, dm [][]float64) float64 {
	cost := 0.0
	for i := 1; i < len(route); i++ {
		cost += dm[route[i-1]][route[i]]
	}

	if len(route)-1 < cluster_num[0] {
		cost += penalty
	}

	return cost
}

func isFeasible(routes [][]int) bool {
	for _, r := range routes {
		//		if len(r)-1 < cluster_num[0] || len(r)-1 > cluster_num[1] {
		if len(r)-1 > cluster_num[1] {

			return false
		}
	}
	return true
}

func copySlice(target, source [][]int) {
	for i, r := range source {
		if len(target[i]) < len(r) {
			target[i] = append(target[i], make([]int, len(r)-len(target[i]))...)
		}
		target[i] = target[i][:len(r)]
		for j, k := range r {
			target[i][j] = k
		}
	}
}
