package main

import (
	"fmt"
	"lml/rand"
	"math"
)

func sa(initSolution [][]int, tm [][]float64, cooling float64) [][]int {
	T := 1.0
	Tmin := 0.001

	if cooling < 0 || cooling >= 1 {
		cooling = 0.99
	}

	var (
		oldFitness  float64
		newFitness  float64
		oldOvertime float64
		newOvertime float64
	)

	oldFitness, oldOvertime = fitness(initSolution, tm)

	var newSolution, currSolution [][]int
	currSolution = make([][]int, len(initSolution))
	for i, r := range initSolution {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}

	for T > Tmin {
		for i := 0; i < 50; i++ {
			dice := rand.Float64()

			newSolution = make([][]int, len(currSolution))
			for i, r := range currSolution {
				newSolution[i] = make([]int, len(r))
				copy(newSolution[i], r)
			}

			// здесь проблема с двойными нулями и потом с пропадающими точами в мутации
			newSolution = mutate(newSolution)

			newFitness, newOvertime = fitness(newSolution, tm)
			if newOvertime > oldOvertime {
				continue
			}

			if newFitness < oldFitness {
				currSolution = newSolution

				ptsOrderMap := map[int]struct{}{}
				for _, r := range initSolution {
					for _, p := range r {
						ptsOrderMap[p] = struct{}{}
					}
				}

				for _, route := range currSolution {
					for _, p := range route {
						if _, ok := ptsOrderMap[p]; ok {
							delete(ptsOrderMap, p)
						}
					}
				}

				if len(ptsOrderMap) > 0 {
					fmt.Println(ptsOrderMap)
					panic("sa: not all points in solution!")
				}

				oldFitness = newFitness
				oldOvertime = newOvertime
				//				fmt.Println("fitness", oldFitness, "routes", len(currSolution))
			} else if dice <= getAcceptanceCoeff(T, oldFitness, newFitness) {

				currSolution = newSolution

				ptsOrderMap := map[int]struct{}{}
				for _, r := range initSolution {
					for _, p := range r {
						ptsOrderMap[p] = struct{}{}
					}
				}

				for _, route := range currSolution {
					for _, p := range route {
						if _, ok := ptsOrderMap[p]; ok {
							delete(ptsOrderMap, p)
						}
					}
				}

				if len(ptsOrderMap) > 0 {
					fmt.Println(ptsOrderMap)
					panic("sa accepted: not all points in solution!")
				}

				oldFitness = newFitness
				oldOvertime = newOvertime

				//				fmt.Println("fitness", oldFitness, "routes", len(currSolution))
			}
		}

		T *= cooling
	}

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
		copy(target[i], r)
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
