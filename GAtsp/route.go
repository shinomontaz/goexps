package solver

import (
	"math"
	"math/rand"
	"time"
)

type Route struct {
	Way     []int // ordered indexes of a cities
	fitness float64
	dMatrix [][]float64
}

func (r *Route) updateStartPoint() {

	var idxFirst int
	for idx, v := range r.Way {
		if v == 0 {
			idxFirst = idx
			break
		}
	}

	r.Way[0], r.Way[idxFirst] = r.Way[idxFirst], r.Way[0]
}

func newRoute(dMatrix [][]float64) *Route {
	rand.Seed(time.Now().UTC().UnixNano())
	route := &Route{Way: rand.Perm(len(dMatrix)), dMatrix: dMatrix}
	route.updateStartPoint()
	return route
}

/*
func doCreate(dMatrix [][]float64) *Route {
	rand.Seed(time.Now().UTC().UnixNano())
	route := &Route{Way: rand.Perm(len(dMatrix)), dMatrix: dMatrix}
	route.updateStartPoint()
	return route
}
*/
func (r *Route) doCrossover(spouse *Route) *Route {
	rand.Seed(time.Now().UnixNano())
	randIndex1 := rand.Intn(len(r.Way) - 1)
	randIndex2 := randIndex1 + rand.Intn(len(r.Way)-randIndex1)

	child := &Route{dMatrix: r.dMatrix}

	// now in child.Way находятся все элементы в тойже последовательности что и в r с пропусками в местах элементов, что попали в интервал [index1, index2] для второго партнера
	//	copy(child.Way, r.Way[randIndex1:randIndex2+1])

	child.Way = append(child.Way, r.Way[randIndex1:randIndex2+1]...)

	for _, i := range spouse.Way {
		if !contains(child.Way, i) {
			child.Way = append(child.Way, i)
		}
	}

	child.updateStartPoint()

	return child
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (r *Route) doMutation(rate float64) {
	rand.Seed(time.Now().UnixNano())
	dice := rand.Float64()

	if dice > rate { // chance to mutate is gone
		return
	}

	randIndex1, randIndex2 := rand.Intn(len(r.Way)-1), rand.Intn(len(r.Way)-1)
	r.Way[randIndex1], r.Way[randIndex2] = r.Way[randIndex2], r.Way[randIndex1]

	r.updateStartPoint()

	_ = r.getFitness()
}

func (r *Route) getFitness() float64 {
	if r.fitness > 0 {
		return r.fitness
	}

	// calculate fitness
	for index, v := range r.Way {
		if (index + 1) < len(r.Way) {
			r.fitness += r.dMatrix[v][r.Way[index+1]]
		}
	}

	r.fitness = math.Pow(r.fitness, 10)
	r.fitness = 1 / (r.fitness + 1)

	return r.fitness
}

func (r *Route) getFitness2() float64 {
	// calculate fitness
	fitness := 0.0
	for index, v := range r.Way {
		if (index + 1) < len(r.Way) {
			fitness += r.dMatrix[v][r.Way[index+1]]
		}
	}

	return fitness
}
