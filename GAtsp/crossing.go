package solver

import "fmt"

func findCrossing(points []*LatLng, route []int) (crossedPairs [][2][2]int) {
	for i := 0; i < len(route)-1; i++ {
		p1 := points[route[i]]
		p2 := points[route[i+1]]

		for j := i + 2; j < len(route)-1; j++ {
			q1 := points[route[j]]
			q2 := points[route[j+1]]

			fmt.Println("checking", [][2]int{[2]int{route[i], route[i+1]}, [2]int{route[j], route[j+1]}})

			if isCrossing(p1, p2, q1, q2) {
				crossedPairs = append(crossedPairs, [2][2]int{[2]int{i, i + 1}, [2]int{j, j + 1}})
			}

		}
	}

	return crossedPairs
}

func isCrossing1(p1 *LatLng, p2 *LatLng, q1 *LatLng, q2 *LatLng) bool {
	//уравнение:
	// x11 + (x12-x11) * t = x21 + (x22 - x21) * t
	// y11 + (y12 - y11) * s = y21 + (y22 - y21) * s
	// определяем s и t => если они оба в интервале [0, 1) значит есть пересечение

	t := (p1.Lat - q1.Lat) / (q2.Lat - q1.Lat - p2.Lat - p1.Lat)
	s := (p1.Lng - q1.Lng) / (q2.Lng - q1.Lng - p2.Lng - p1.Lng)

	return t > 0.0 && t < 1.0 && s > 0.0 && s < 1.0
}

func isCrossing2(p1 *LatLng, p2 *LatLng, q1 *LatLng, q2 *LatLng) bool {
	//уравнение:
	// x - x2 / x2-x1 = y - y2 / y2 - y1
	// y =
	x1 := p1.Lat
	x2 := p2.Lat
	x3 := q1.Lat
	x4 := q2.Lat
	y1 := p1.Lng
	y2 := p2.Lng
	y3 := q1.Lng
	y4 := q2.Lng

	xCrossing := (x4*(y4-y3)/(x4-x3) - y3 + y2 - x2*(y2-y1)/(x2-x1)) / ((y4-y3)/(x4-x3) - (y2-y1)/(x2-x1))

	return ((xCrossing > x1 && xCrossing < x2) || (xCrossing > x2 && xCrossing < x1)) && ((xCrossing > x3 && xCrossing < x4) || (xCrossing > x4 && xCrossing < x3))
}

func isCrossing(p1 *LatLng, p2 *LatLng, q1 *LatLng, q2 *LatLng) bool {
	v1 := vectorProduct(q2.Lat-q1.Lat, q2.Lng-q1.Lng, p1.Lat-q1.Lat, p1.Lng-q1.Lng)
	v2 := vectorProduct(q2.Lat-q1.Lat, q2.Lng-q1.Lng, p2.Lat-q1.Lat, p2.Lng-q1.Lng)
	v3 := vectorProduct(p2.Lat-p1.Lat, p2.Lng-p1.Lng, q1.Lat-p1.Lat, q1.Lng-p1.Lng)
	v4 := vectorProduct(p2.Lat-p1.Lat, p2.Lng-p1.Lng, q2.Lat-p1.Lat, q2.Lng-p1.Lng)

	return v1*v2 < 0 && v3*v4 < 0
}

func vectorProduct(a, b, c, d float64) float64 {
	return a*d - c*b
}
