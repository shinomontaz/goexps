package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Point struct {
	x float64
	y float64
}

func initAll() {
	//	rand.Seed(time.Now().UTC().UnixNano())
	rand.Seed(time.Now().UTC().UnixNano())
}

func generatePoints(n int) []Point {

	limitX := 100.0
	limitY := 100.0

	points := make([]Point, 0, n)
	for i := 0; i < n; i++ {
		x := rand.Float64() * limitX
		y := rand.Float64() * limitY
		points = append(points, Point{x: x, y: y})
	}

	return points
}

func calcDistances(points []Point) [][]float64 {
	dMatrix := make([][]float64, 0, len(points))
	for _, from := range points {
		row := make([]float64, 0, len(points))
		for _, to := range points {
			row = append(row, distance(from, to))
		}
		dMatrix = append(dMatrix, row)
	}

	return dMatrix
}

func distance(p, q Point) float64 {
	return math.Sqrt(math.Pow((p.x-q.x), 2) + math.Pow(p.y-q.y, 2))
}

type Route []int

var Population []Route

func main() {
	initAll()
	N := 10
	points := generatePoints(N)
	//	dMatrix := calcDistances(points)

	pop := GAgeneratePopulation(N*100, points)

	var currBest Route
	var everBest Route

	sameCount := 0
	sameLimit := 1000

	everBest = pop[0]
	fmt.Println("Started: ", everBest, " - ", everBest.Score(points))
	for i := 0; i < 100000; i++ {
		currBest = GAfindBest(pop, points)
		if everBest.Score(points) < currBest.Score(points) {
			everBest = currBest
			sameCount = 0
		} else {
			sameCount++
		}

		if sameCount > sameLimit {
			break
		}
		pop = GAevolvePopulation(pop)

		fmt.Println("In progress: ", currBest.Score(points), " - ", everBest.Score(points), " - ", sameCount)

	}

	fmt.Println(everBest, " - ", everBest.Score(points))
}

func NewRoute(points []Point) Route {
	return rand.Perm(len(points))
}

func (r Route) Score(points []Point) float64 {
	sum := 0.0
	for _, i := range r {
		if i < len(points)-1 {
			sum += distance(points[i], points[i+1])
		}
	}

	return 1 / (sum + 1)
}

func (r Route) Mutate() Route {
	randIndex1, randIndex2 := rand.Intn(len(r)), rand.Intn(len(r))
	r[randIndex1], r[randIndex2] = r[randIndex2], r[randIndex1]
	return r
}

func GAgeneratePopulation(m int, points []Point) []Route {
	population := make([]Route, 0, m)
	for i := 0; i < m; i++ {
		population = append(population, NewRoute(points))
	}

	return population
}

func GAevolvePopulation(population []Route) []Route {
	// create new population in place of old one

	// just mutate all
	newPopulation := make([]Route, 0, len(population))

	for i := 0; i < len(population); i++ {

		parent1 := GAPick(population)
		parent2 := GAPick(population)

		child := GACrossover(parent1, parent2)

		newPopulation = append(newPopulation, child.Mutate())
	}

	return newPopulation
}

func GAPick(population []Route) Route {
	randIndex := rand.Intn(len(population))

	selected := make(Route, 0, len(population[randIndex]))

	selected = append(selected, population[randIndex]...)

	return selected
}

func GACrossover(parent1, parent2 Route) Route {
	child := make([]int, 0, len(parent1))                         //                               <-------------> = len(parent1) - randIndex1
	randIndex1 := rand.Intn(len(parent1) - 1)                     // secondcan be bigger _________(.)_____________
	randIndex2 := randIndex1 + rand.Intn(len(parent1)-randIndex1) // 					 ___________________(.)___

	child = append(child, parent1[randIndex1:randIndex2+1]...)

	for _, i := range parent2 {
		if Route(child).IndexOf(i) < 0 {
			child = append(child, i)
		}
	}

	return child
}

func GAfindBest(population []Route, points []Point) Route {
	sort.Slice(population, func(i, j int) bool {
		return population[i].Score(points) > population[j].Score(points)
	})

	return population[0]
}

func (r Route) IndexOf(i int) int {
	for idx, v := range r {
		if v == i {
			return idx
		}
	}
	return -1
}
