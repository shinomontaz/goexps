package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"lml/rand"
	"os"
	"strconv"

	"log"
)

var (
	OsrmURL       string
	courier_shift float64
	box_time      float64
)

func init() {
	courier_shift = 7 * 60 * 60            // 	MAX_HOURS = 7
	box_time = 15 * 60                     // Время на разгрузку в ПВЗ - 15 минут
	OsrmURL = "http://91.228.153.227:5000" // "http://osrm-gazelle.wbdispatch.svc.k8s.stage-dp"
	rand.Init(1 << 24)
}

func main() {
	wh := LatLong{
		Lat: 55.386871,
		Lng: 37.588898,
	}

	ptsMap := map[LatLong]int{}
	pts := []LatLong{wh}

	f, err := os.Open("points.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	_, err = csvReader.Read() // read header
	if err != nil {
		log.Fatal(err)
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		lat, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		lng, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		ll := LatLong{
			Lat: lat,
			Lng: lng,
		}
		pts = append(pts, ll)

		ptsMap[ll] += 1
	}

	non_unique := map[LatLong]int{}
	for p, num := range ptsMap {
		if num > 1 {
			non_unique[p] = num
		}
	}

	fmt.Println("non_unique: ", non_unique)

	// матрицы расстояний
	l := len(pts)
	dm := make([][]float64, l)
	tm := make([][]float64, l)

	for i := 0; i < l; i++ {
		dm[i] = make([]float64, l)
		tm[i] = make([]float64, l)
	}

	max_len := 1000
	chunks := int(float64(l) / float64(max_len))
	if l%max_len > 0 {
		chunks++
	}

	var src, dst []int
	for j := 0; j < chunks; j++ {
		src = make([]int, 0, l)
		for id := j * max_len; id < l && id < (j+1)*max_len; id++ {
			src = append(src, id)
		}
		fmt.Println("start receiving ", j, " out of ", chunks)
		for jj := 0; jj < chunks; jj++ {
			fmt.Println("start receiving ", jj, " out of ", chunks, " for row ", j)

			dst = make([]int, 0, l)
			for id := jj * max_len; id < l && id < (jj+1)*max_len; id++ {
				dst = append(dst, id)
			}
			ttm, _ := GetMatrices2(pts, src, dst) // times and distances
			for iii, row := range ttm {
				for jjj, p := range row {
					tm[iii+j*max_len][jjj+jj*max_len] = p * 0.8
				}
			}
		}
	}

	fmt.Println(len(dm))

	//	initialRoutes := initialFlower(pts, tm, dm)
	initialRoutes := initialGreedy(pts, tm)
	//	initialRoutes := initialRandomGreedy(pts, tm)
	fmt.Println(len(initialRoutes))

	currSolution := make([][]int, len(initialRoutes))
	for i, r := range initialRoutes {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}
	var (
		newFitness   float64
		newOvertime  float64
		oldFitness   float64
		oldOvertime  float64
		newSolution  [][]int
		copySolution [][]int
	)
	oldFitness, oldOvertime = fitness(initialRoutes, tm)
	fitnesTreshold := oldFitness

	is_splitted := false
	for i := 0; i < 100; i++ {
		fmt.Println("iteration:", i)

		if i == 20 && oldOvertime > courier_shift && !is_splitted {
			// split routes
			num := int(oldOvertime / courier_shift)
			currSolution = splitMaxRoute(currSolution, num, tm)
			oldFitness, oldOvertime = fitness(currSolution, tm)
			is_splitted = true
		}

		if oldOvertime < 2*courier_shift && !is_splitted {
			copySolution = make([][]int, len(currSolution))
			for i, r := range currSolution {
				copySolution[i] = make([]int, len(r))
				copy(copySolution[i], r)
			}
			newSolution = eliminateRandom(copySolution, tm)
			newSolution = sa(newSolution, tm, 0.99)

			newFitness, newOvertime = fitness(newSolution, tm)
			if newOvertime > oldOvertime {
				continue
			}

			if newFitness < oldFitness || newFitness < fitnesTreshold {
				currSolution = newSolution
				oldFitness = newFitness
				oldOvertime = newOvertime
			}
		} else {
			currSolution = sa(currSolution, tm, 0.99)
			newFitness, newOvertime = fitness(currSolution, tm)
			oldFitness = newFitness
			oldOvertime = newOvertime
		}

		fmt.Println("main cycle update:", len(currSolution), oldFitness, oldOvertime)
	}

	finalRoutes := sa(currSolution, tm, 0.9999)
	newFitness, newOvertime = fitness(finalRoutes, tm)
	max, mean, num := max_overtime(finalRoutes, tm)

	fmt.Println(len(finalRoutes), newFitness, newOvertime, max, mean, num)
}
