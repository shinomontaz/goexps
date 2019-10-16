package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type LatLng struct {
	Lat float64
	Lng float64
}

func main() {
	rand.Seed(time.Now().UnixNano())
	points := createPoints(100)
	dMatrix := calcDistances(points)

	route := Greedy(points, dMatrix)
	//	route := rand.Perm(len(dMatrix))
	//	updateStartPoint(route)

	fmt.Println("initial: ", route, " - ", getFitness(route, dMatrix))

	annealed := annealing(dMatrix, route)

	fmt.Println("annealed: ", annealed, " - ", getFitness(annealed, dMatrix))
	drawSolution("SA2.png", points, annealed)
}

func annealing(dMatrix [][]float64, currSolution []int) []int {
	T := 1.0
	Tmin := 0.001
	cooling := 0.999
	oldEnergy := getFitness(currSolution, dMatrix)

	for T > Tmin {
		newSolution := mutate(currSolution)
		newEnergy := getFitness(newSolution, dMatrix)
		if newEnergy < oldEnergy {
			currSolution = newSolution
			oldEnergy = newEnergy

		} else {
			dice := rand.Float64()
			if dice > getAcceptanceCoeff(T, oldEnergy, newEnergy) {
				currSolution = newSolution
				oldEnergy = newEnergy
			}
		}
		T *= cooling
	}
	return currSolution
}

func getAcceptanceCoeff(T float64, oldEnergy, newEnergy float64) float64 {
	return math.Exp((newEnergy - oldEnergy) / T)
}

func getFitness(route []int, dMatrix [][]float64) (fitness float64) {
	// calculate fitness
	// for i := 0; i < len(route)-1; i++ {
	// 	fitness += dMatrix[route[i]][route[i+1]]
	// }
	for index, v := range route {
		if (index + 1) < len(route) {
			fitness += dMatrix[v][route[index+1]]
		}
	}

	return fitness //1 / (fitness + 1)
}

func mutate(route []int) (mutated []int) {
	mutated = append(mutated, route...)
	randIndex1, randIndex2 := 1+rand.Intn(len(route)-1), 1+rand.Intn(len(route)-1)
	mutated[randIndex1], mutated[randIndex2] = mutated[randIndex2], mutated[randIndex1]

	//	updateStartPoint(mutated)

	return mutated
}

func updateStartPoint(route []int) {
	var idxFirst int
	for idx, v := range route {
		if v == 0 {
			idxFirst = idx
			break
		}
	}

	route[0], route[idxFirst] = route[idxFirst], route[0]
}

func createPoints(n int) []*LatLng {
	res := make([]*LatLng, 0)
	for i := 0; i < n; i++ {
		res = append(res, &LatLng{
			Lat: rand.Float64() * 100,
			Lng: rand.Float64() * 100,
		})
	}
	return res
}

func calcDistances(points []*LatLng) [][]float64 {
	res := make([][]float64, 0)
	for _, from := range points {
		row := make([]float64, 0)
		for _, to := range points {
			row = append(row, getDistance(from, to))
		}
		res = append(res, row)
	}

	return res
}

func getDistance(from, to *LatLng) float64 {
	if from == to {
		return 0
	}
	return math.Sqrt(math.Pow(from.Lat-to.Lat, 2) + math.Pow(from.Lng-to.Lng, 2))
}

func drawSolution(filename string, points []*LatLng, route []int) {
	length := 500.0
	myImage := image.NewRGBA(image.Rect(0, 0, int(length), int(length)))
	outputFile, _ := os.Create(filename)
	defer outputFile.Close()
	myred := color.RGBA{200, 0, 0, 255}

	for i := 0; i < len(points); i++ {
		x1 := length * points[i].Lat / 100.0
		y1 := length * points[i].Lng / 100.0

		redRect := image.Rect(int(x1-2), int(y1-2), int(x1+2), int(y1+2))

		// create a red rectangle atop the green surface
		draw.Draw(myImage, redRect, &image.Uniform{myred}, image.ZP, draw.Src)

		addLabel(myImage, int(x1+2), int(y1-2), fmt.Sprintf("%d", i))
	}

	for j := 0; j < len(points)-1; j++ {
		i := route[j]
		k := route[j+1]
		x1 := length * points[i].Lat / 100.0
		y1 := length * points[i].Lng / 100.0

		x2 := length * points[k].Lat / 100.0
		y2 := length * points[k].Lng / 100.0

		addLine(myImage, x1, y1, x2, y2)
	}

	png.Encode(outputFile, myImage)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func addLine(img *image.RGBA, x1, y1, x2, y2 float64) {
	y := y1

	start := math.Min(x1, x2)
	end := math.Max(x1, x2)

	for x := int(start); x < int(end); x++ {
		y = (float64(x)-x1)*(y2-y1)/(x2-x1) + y1
		img.Set(x, int(y), color.RGBA{200, 0, 0, 255})
	}
}

func Greedy(points []*LatLng, dMatrix [][]float64) []int {
	current := 0
	visited := make(map[int]int)
	visited[0] = -1
	last := -1
	for {
		minDist := math.MaxFloat64
		minK := -1
		for k := range points {
			if _, ok := visited[k]; ok {
				continue
			}
			if dMatrix[current][k] < minDist {
				minDist = dMatrix[current][k]
				minK = k
			}
		}
		if minK == -1 {
			break
		} else {
			visited[minK] = current
			current = minK
			last = current
		}

	}
	reversedWay := make([]int, 0)
	current = last
	for current != -1 {
		reversedWay = append(reversedWay, current)
		current = visited[current]

	}
	way := make([]int, 0, len(reversedWay))
	for i := len(reversedWay) - 1; i >= 0; i-- {
		way = append(way, reversedWay[i])
	}

	return way
}
