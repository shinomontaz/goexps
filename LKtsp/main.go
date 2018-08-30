package main

import (
	"fmt"
	"math/rand"
	"time"
)

type LatLng struct {
	Lat float64
	Lng float64
}

type PathFactory struct {
	N       int
	dMatrix [][]float64
	points  []*LatLng
}

func (f *PathFactory) Create() Path {
	path := Path{way: rand.Perm(f.N), points: f.points, neighbours: make(map[int][]int, f.N)}

	for _, i := range path.way {
		path.neighbours[i] = make([]int, 0)

		for j := 0; j < f.N; j++ {
			if f.dMatrix[i][j] > 0 {
				path.neighbours[i] = append(path.neighbours[i], j)
			}
		}
	}

	return path
}

func main() {

	improved := true
	var N = 10
	rand.Seed(time.Now().UnixNano())

	points := createPoints(N)
	dMatrix := calcDistances(points)

	PFactory := &PathFactory{N: N, dMatrix: dMatrix, points: points}
	path := PFactory.Create()

	fmt.Println(path)

	panic("!")

	solutions := make([][]int, 0, 0)

	for improved != false {
		improved = improve(&path)
		solutions = append(solutions, path.way)
	}
}

func improve(path *Path) bool {
	//	Find all valid 2-opt moves and try them
	for _, t1 := range path.way {
		around := path.Around(t1)

		for t2 := range around {
			broken := []int{t1, t2} //set([makePair(t1, t2)])
			// Initial savings
			gain := getDistance(path.points[t1], path.points[t2]) // TSP.dist(t1, t2)

			close := path.closest(t2, tour, gain, broken, set())

			// Number of neighbours to try
			tries := 5
			for {
				//			for t3, (_, Gi) in close {
				// Make sure that the new node is none of t_1's neighbours
				// so it does not belong to the tour.
				isIn := false
				for _, it := range around {
					if it == t3 {
						isIn = true
						break
					}
				}
				if isIn {
					continue
				}

				joined := []int{t2, t3} // set([makePair(t2, t3)])

				// The positive Gi is taken care of by `closest()`

				if path.chooseX(tour, t1, t3, Gi, broken, joined) {
					// Return to Step 2, that is the initial loop
					return true
				}
				// Else try the other options

				tries--
				// Explored enough nodes, change t_2
				if tries == 0 {
					break
				}
			}
		}

	}

	return false
}

func (path *Path) closest() {

}
