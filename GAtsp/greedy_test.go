package solver

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func getDistance(from, to *LatLng) float64 {
	if from == to {
		return 0
	}
	return math.Sqrt(math.Pow(from.Lat-to.Lat, 2) + math.Pow(from.Lng-to.Lng, 2))
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

func TestGreedyTwoPoints(t *testing.T) {

	points := createPoints(5)

	distances := calcDistances(points)
	//	grRes := greedy(points, distances)
	genRes := New(points, distances).Solve(0.01, 1000, 100000)

	fmt.Println(genRes.Way, " - ", genRes.getFitness2())

	drawSolution("before.png", points, genRes)

	res := findCrossing(points, genRes.Way)

	for _, crossPair := range res {
		fmt.Println("fixing: ", crossPair, " - ", genRes.Way[crossPair[0][0]], "reverting ", crossPair[0][1], " - ", crossPair[1][0])
		reverseSlice(genRes.Way[crossPair[0][1] : crossPair[1][0]+1])
		fmt.Println("fixed: ", genRes.Way, " - ", genRes.getFitness2())
	}

	fmt.Println(genRes.Way, " - ", genRes.getFitness2())

	drawSolution("after.png", points, genRes)

}

func drawSolution(filename string, points []*LatLng, genRes *Route) {
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
		i := genRes.Way[j]
		k := genRes.Way[j+1]
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

func reverseSlice(input []int) []int {
	for i := len(input)/2 - 1; i >= 0; i-- {
		opp := len(input) - 1 - i
		input[i], input[opp] = input[opp], input[i]
	}
	return input
}
