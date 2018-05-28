package solver

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type LatLng struct {
	Lat float64
	Lng float64
}

type City struct {
	LatLng
	distances []float64 // distances to other cities - it is quite specific param only for this TSP setup
}

type Gsolver struct {
	population []*Route
	// PopSize       int
	// chanceMutate float64
	// limitLocalMax int
	// maxGenerations int
	cities       []*City
	points       []*LatLng
	popSize      int
	totalFitness float64 // sum of fitnesses of current population
	dMatrix      [][]float64
}

// New function returns instance of a solver with initialized fields
func New(pointList []*LatLng, dMatrix [][]float64) *Gsolver {
	gSolver := &Gsolver{points: pointList}
	gSolver.init(dMatrix)
	gSolver.dMatrix = dMatrix
	return gSolver
}

// init function
// make query to get distance matrix
// initalizes slice of cities with proper distance list for every one - taking a matrix row
func (s *Gsolver) init(dMatrix [][]float64) {
	n := len(s.points)
	// prepare distance matrix
	s.cities = make([]*City, 0, n)
	for index, point := range s.points {
		city := &City{LatLng: *point, distances: make([]float64, 0, n)}
		city.distances = dMatrix[index]
		s.cities = append(s.cities, city)
	}

	s.popSize = len(s.cities) * 20
}

// main method of GA:
// create initial random population
// loop throw generations
// estimate each population and find best individual
// crossover by fitness score, mutate and continue
func (s *Gsolver) Solve(mutationRate float64, stopAfter int, maxGenerations int) *Route {

	start := time.Now().UTC()

	var bestEver *Route
	var bestCurr *Route
	lastsRecord := 0
	s.createPopulation()

	fmt.Println("initial population: ", time.Since(start))

	for i := 0; i < maxGenerations; i++ {
		bestCurr = s.getBest()
		if bestEver == nil || bestCurr.getFitness() > bestEver.getFitness() {
			bestEver = bestCurr
			lastsRecord = 0
		} else {
			lastsRecord++
		}

		if lastsRecord > stopAfter {
			break
		}
		s.evolvePopulation()
	}

	return bestEver
}

func (s *Gsolver) createPopulation() {
	s.population = make([]*Route, 0, s.popSize)
	s.totalFitness = 0

	var wg sync.WaitGroup
	chRoute := make(chan *Route)
	for i := 0; i < s.popSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chRoute <- newRoute(s)
		}()
	}

	go func() {
		wg.Wait()
		close(chRoute)
	}()

	for route := range chRoute {
		s.population = append(s.population, route)
		s.totalFitness += route.getFitness()
	}
}

func (s *Gsolver) evolvePopulation() {
	newTotalFitness := 0.0
	newPopulation := make([]*Route, 0, s.popSize)

	var wg sync.WaitGroup
	chChilds := make(chan *Route, s.popSize)
	for i := 0; i < s.popSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			parent1 := s.pick()
			parent2 := s.pick()
			child := parent1.doCrossover(parent2)
			child.doMutation(0.01)
			chChilds <- child
		}()
	}

	go func() {
		wg.Wait()
		close(chChilds)
	}()

	for child := range chChilds {
		newPopulation = append(newPopulation, child)
		newTotalFitness += child.getFitness()
	}

	s.population = newPopulation
	s.totalFitness = newTotalFitness
}

func (s *Gsolver) pick() *Route {
	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Float64()

	i := 0
	for ; random > 0; i++ {
		random -= (s.population[i].getFitness() / s.totalFitness)
	}
	i--

	return &Route{context: s, Way: s.population[i].Way, fitness: s.population[i].fitness}
}

func (s *Gsolver) getBest() *Route {
	sort.Slice(s.population, func(i, j int) bool {
		return s.population[i].getFitness() > s.population[j].getFitness()
	})
	return s.population[0]
}
