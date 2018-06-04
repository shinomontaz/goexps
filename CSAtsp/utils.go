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

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func createPoints(n int) []*Point {
	res := make([]*Point, 0)
	for i := 0; i < n; i++ {
		startSec := MegaNow + int64(rand.Intn(17)*30*60)
		res = append(res, &Point{
			Lat:   rand.Float64() * 100,
			Lng:   rand.Float64() * 100,
			Start: startSec,
			End:   startSec + 30*60,
		})
	}

	return res
}

func calcDistances(points []*Point) [][]float64 {
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

func getDistance(from, to *Point) float64 {
	if from == to {
		return 0
	}
	return math.Sqrt(math.Pow(from.Lat-to.Lat, 2) + math.Pow(from.Lng-to.Lng, 2))
}

func drawSolution(filename string, points []*Point, route []int) {
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

func Greedy() []int {
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
