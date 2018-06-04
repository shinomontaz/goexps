package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Point struct {
	Lat   float64
	Lng   float64
	Start int64 // timestamp
	End   int64
}

var dMatrix [][]float64
var points []*Point
var MegaNow int64

func init() {

	rand.Seed(time.Now().UnixNano())
	today := time.Now()
	MegaNow = time.Date(today.Year(), today.Month(), today.Day(), 8, 0, 0, 0, time.UTC).Unix()
	points = createPoints(10)
	dMatrix = calcDistances(points)
}

func main() {

	route := Greedy()

	csa := Csa(route)

	PrintPoints()

	fmt.Println("initial: ", route, " - ", _fu(route), " sum of out of TW: ", _ro(csa))
	PrintRoute(route)
	//	drawSolution("CSA1.png", points, route)
	fmt.Println("csa: ", csa, " - ", _fu(csa), " sum of out of TW: ", _ro(csa))
	PrintRoute(csa)

	//	drawSolution("CSA2.png", points, csa)
}

func Csa(currSolution []int) []int {
	T := 1.0
	Tmin := 0.0001
	cooling := 0.999
	P := 1000000.0
	releasing := 0.999
	Pmin := 0.0001

	oldEnergy := _nu(currSolution, P)

	for T > Tmin && P > Pmin {
		newSolution := mutate(currSolution)
		newEnergy := _nu(newSolution, P)
		if newEnergy < oldEnergy {
			currSolution = newSolution
			oldEnergy = newEnergy

		} else {
			dice := rand.Float64()
			if dice > getAcceptanceCoeff(T, P, oldEnergy, newEnergy) {
				currSolution = newSolution
				oldEnergy = newEnergy
			}
		}
		T *= cooling
		P *= releasing
	}
	return currSolution
}

func getAcceptanceCoeff(T float64, P float64, oldEnergy, newEnergy float64) float64 {
	return math.Exp((newEnergy - oldEnergy) / T)
}

func _nu(route []int, P float64) float64 {
	return float64(_fu(route)) + P*float64(_ro(route))
}

func _c(i, j int) int64 {
	return int64(dMatrix[i][j] / 13.0)
}

func _ai(index int, route []int) int64 {
	if index == 0 {
		return MegaNow
	}
	return _di(index-1, route) + _c(index-1, index)
}

func _di(index int, route []int) int64 {
	if index == 0 {
		return MegaNow
	}
	return Max(_ai(index, route), points[route[index]].Start)
}

func _ro(route []int) (res int64) {
	for index := range route {
		res += Max(0, _di(index, route)-points[index].End)
	}
	return res
}

func _fu(route []int) (res int64) {
	for index, v := range route {
		if (index + 1) < len(route) {
			res += _c(v, route[index+1])
		}
	}

	return res
}

func mutate(route []int) (mutated []int) {
	mutated = append(mutated, route...)

	randIndex1, randIndex2 := 1+rand.Intn(len(route)-1), 1+rand.Intn(len(route)-1)
	mutated[randIndex1], mutated[randIndex2] = mutated[randIndex2], mutated[randIndex1]

	return mutated
}

func PrintRoute(route []int) {
	for _, v := range route {
		fmt.Println(v, " [arrival: ", time.Unix(_ai(v, route), 0).UTC(), ", leaving: ", time.Unix(_di(v, route), 0).UTC())
	}
}

func PrintPoints() {
	for _, v := range points {
		fmt.Println(v, " [start: ", time.Unix(v.Start, 0).UTC(), ", end: ", time.Unix(v.End, 0).UTC())
	}
}
