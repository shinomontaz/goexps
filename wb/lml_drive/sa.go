package main

import (
	"fmt"
	"lml-drive/rand"
	"lml-drive/types"
	"math"
)

func sa(initSolution [][]int, pts []types.Point, dm [][]float64, tm [][]float64, cooling float64) [][]int {
	T := 1.0
	Tmin := 0.001

	if cooling < 0 || cooling >= 1 {
		cooling = 0.99
	}

	var (
		oldFitness    float64
		newFitness    float64
		fitTreshold   float64
		newMaxOvershk float64
		oldOvershk    int
		oldUndershk   int
		newOvershk    int
		newUndershk   int
	)

	oldFitness, oldOvershk, oldUndershk, _ = fitness(initSolution, pts, dm)

	var newSolution, currSolution [][]int
	currSolution = make([][]int, len(initSolution))
	for i, r := range initSolution {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}

	newSolution = make([][]int, len(currSolution))

	for T > Tmin {
		fitTreshold = oldFitness
		for i := 0; i < 1000; i++ {

			newSolution = newSolution[:len(currSolution)]
			copySlice(newSolution, currSolution)

			dice := rand.Float64()

			if dice < 0.01 {
				// remove 1 route
				newSolution = educate(newSolution, pts, dm, 1.0 /*0.9*/)
			} else {
				newSolution = mutate(newSolution)
			}

			newFitness, newOvershk, newUndershk, newMaxOvershk = fitness(newSolution, pts, dm)

			if newUndershk > oldUndershk || newOvershk > oldOvershk || isFeasible(newSolution, tm) > 0 {
				continue
			}

			// if ((newUndershk > oldUndershk || newOvershk > oldOvershk) && len(newSolution) >= len(currSolution)) || newMaxOvershk > 0.05 || isFeasible(newSolution, tm) > 0 {
			// 	continue
			// }

			if newFitness < oldFitness {

				//				currSolution = newSolution
				currSolution = currSolution[:len(newSolution)]
				copySlice(currSolution, newSolution)
				oldFitness = newFitness
				oldOvershk = newOvershk
				oldUndershk = newUndershk

				//				fmt.Println("fitness", oldFitness, newOvershk, "routes", len(currSolution))
				if newMaxOvershk > 0 {
					fmt.Println(newMaxOvershk)
				}
			} else if dice <= getAcceptanceCoeff(T, oldFitness, newFitness) && newFitness < fitTreshold {
				//				currSolution = newSolution
				currSolution = currSolution[:len(newSolution)]
				copySlice(currSolution, newSolution)

				oldFitness = newFitness
				oldOvershk = newOvershk
				oldUndershk = newUndershk

				//				fmt.Println("accepted: ", oldFitness, newOvershk, "routes", len(currSolution))
			}
		}

		T *= cooling
	}
	fmt.Println("annealed: ", oldFitness, oldUndershk, oldOvershk, "routes", len(currSolution))

	newSolution = make([][]int, len(currSolution))
	j := 0
	for _, r := range currSolution {
		newSolution[j] = make([]int, len(r))
		if len(r) > 1 {
			copy(newSolution[j], r)
			j++
		}
	}

	newSolution = newSolution[:j]

	ptsOrderMap := map[int]struct{}{}
	for _, r := range initSolution {
		for _, p := range r {
			ptsOrderMap[p] = struct{}{}
		}
	}

	for _, route := range newSolution {
		for _, p := range route {
			if _, ok := ptsOrderMap[p]; ok {
				delete(ptsOrderMap, p)
			}
		}
	}

	if len(ptsOrderMap) > 0 {
		fmt.Println(ptsOrderMap)
		panic("not all points in solution!")
	}

	return newSolution
}

func getAcceptanceCoeff(T float64, oldEnergy, newEnergy float64) float64 {
	// oldEnergy < newEnergy => степень будет отрицательной. И чем меньше температура, тем меньше знаменатель и больше отрицательная степень
	return math.Exp((oldEnergy - newEnergy) / (oldEnergy * T))
}

func copySlice(target, source [][]int) {
	for i, r := range source {
		if len(target[i]) < len(r) {
			target[i] = append(target[i], make([]int, len(r)-len(target[i]))...)
		}
		target[i] = target[i][:len(r)]
		//		copy(target[i], r)
		for j, k := range r {
			target[i][j] = k
		}
	}
}

func mutate(mutated [][]int) [][]int {
	dice := rand.Float64()
	if dice <= 0.25 {
		// swap move
		return swapmove(mutated)
	} else if dice > 0.25 && dice <= 0.75 {
		//		insertion move
		return insertmove(mutated)
	}
	//	2-opt move
	return twooptmove(mutated)
}
