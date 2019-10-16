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

	fmt.Println("initial: ", route, " - ", getFitness(route, dMatrix))
	drawSolution("Greedy.png", points, route)

}

func FI(points []*LatLng, dm [][]float64) {
	// detect 2 longest and create edge
	// detect 3 far from edge
	// loop throgh others
	// a) Find the farthest city yet to be incorporated, from any point on the mini-tour (city 'd').
	// b) Note the edge to which this city is closest.
	// c) Erase this edge, and create two new edges, between city 'd' each of the cities at the ends of the erased closest edge.

	// find biggest in matrix

	var pointsMap map[int]struct{}
	for i := range points {
		pointsMap[i] = struct{}{}
	}

	max := 0.0
	var maxI, maxJ int
	for i, row := range dm {
		for j, d := range row {
			if d > max {
				maxI = i
				maxJ = j
			}
		}
	}

	edge1 := []*LatLng{points[maxI], points[maxJ]}
	route := []int{maxI, maxJ}
	delete(pointsMap, maxI)
	delete(pointsMap, maxJ)

	// find far point from edge
	max = 0
	maxI = -1
	for i := range pointsMap {
		// расстояние от точки points[i] до первого ребра считаем
	}
	delete(pointsMap, maxI)
	route = append(route, maxI)

	// main loop
	for true {
		if len(pointsMap) == 0 {
			break
		}

		max = 0
		maxI = -1
		for i := range pointsMap {
			for j := range route {
				if dm[i][j] > max {
					max = dm[i][j]
					maxI = i
				}
			}
		}
		delete(pointsMap, maxI)
		route = append(route, maxI)
	}

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
