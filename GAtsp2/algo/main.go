package algo

import "math/rand"

type Individual interface {
	Fitness() float64
	Mutate(rng *rand.Rand)
	Crossover(parner Individual) (child Individual)
	Clone() Individual
}

type IFactory interface {
	Init(rng *rand.Rand)
	Create() Individual
}
