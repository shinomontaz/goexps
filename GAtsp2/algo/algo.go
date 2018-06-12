package algo

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

//type IndividualFactory func(rng *rand.Rand) Individual

type Solver struct {
	If IFactory
	// Required fields
	//	NewIndividual NewIndividual `json:"-"`
	PopSize int `json:"-"` // Number of Individuls per Population

	// Optional fields
	Rnd *rand.Rand `json:"-"`

	// Fields generated at runtime
	Population   []Individual
	Best         Individual `json:"hall_of_fame"` // Sorted best Individuals ever encountered
	totalFitness float64
	Generations  int `json:"generations"` // Number of generations the GA has been evolved
	NotChanged   int
}

func (s *Solver) Initialize() {
	// create initial population
	if s.Rnd == nil {
		s.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if s.Generations == 0 {
		s.Generations = s.PopSize * 10000
	}

	if s.NotChanged == 0 {
		s.NotChanged = 1000
	}

	s.If.Init(s.Rnd)

	s.Population = make([]Individual, 0, s.PopSize)
	for i := 0; i < s.PopSize; i++ {
		s.Population = append(s.Population, s.If.Create())
	}

	s.Best = s.Population[0]

	// Sort it's Individuals
	//	ga.Populations[i].Individuals.SortByFitness()
}

func (s *Solver) pick() Individual {
	return s.Population[0]
}

func (s *Solver) evolve() {
	// create new population in place of current
	// crossover
	// mutate

	newPopulation := make([]Individual, 0, s.PopSize)
	var fitnessSum float64

	var wg sync.WaitGroup
	chEvolved := make(chan Individual, s.PopSize)
	for i := 0; i < s.PopSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parent1 := s.pick()
			//			parent2 := s.pick()
			//			child := parent1.Crossover(parent2)
			child := parent1.Clone()
			child.Mutate(s.Rnd)
			chEvolved <- child
		}()
	}

	go func() {
		wg.Wait()
		close(chEvolved)
	}()

	for newIdividual := range chEvolved {
		newPopulation = append(newPopulation, newIdividual)
		fitnessSum += newIdividual.Fitness()
	}

	s.Population = newPopulation
	s.totalFitness = fitnessSum
}

func (s *Solver) findBest() Individual {
	sort.SliceStable(s.Population, func(i, j int) bool {
		return s.Population[i].Fitness() > s.Population[j].Fitness() // it is a "less" function, so we need bigger first
	})

	return s.Population[0]
}

func (s *Solver) Solve() {
	var currBest Individual
	countNotChanged := 0
	for i := 0; i < s.Generations; i++ {
		currBest = s.findBest()
		if s.Best.Fitness() <= currBest.Fitness() {
			countNotChanged++
			s.Best = currBest
		} else {
			countNotChanged = 0
		}

		fmt.Println("sdfds: ", s.Best.Fitness())

		if countNotChanged > s.NotChanged {
			return
		}
		s.evolve()
	}
}
