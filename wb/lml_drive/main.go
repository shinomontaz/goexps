package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"lml-drive/rand"
	"os"
	"strconv"

	"log"

	"net/http"
	_ "net/http/pprof"
)

type dist_shk_f func(dist float64) (int, int)
type dist_max_f func(dist float64) float64

var (
	OsrmURL string
	fshk    dist_shk_f
	fdist   dist_max_f

	max_pts     int
	office_time float64
)

/*
0-50	6	 800-1000 	 5 000
50-100	5	 1000-1200 	 5 000
100-200	4	 1200-1500 	 4 800
200-300	3	 1500-2000 	 4 200
300-500	2	 2000-2500 	 4 000
500+	1	 2500-3000 	 3 000
*/

func init() {
	OsrmURL = "http://osrm-gazelle.wbdispatch.svc.k8s.stage-dp" // http://91.228.153.227:5000" //
	rand.Init(1 << 24)
	fshk = func(dist float64) (int, int) {
		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 4500, 6250
		case dist > 5.0 && dist <= 100.0:
			return 4500, 6250
		case dist > 100.0 && dist <= 200.0:
			return 4320, 6000
		case dist > 200.0 && dist <= 300.0:
			return 3780, 5250
		case dist > 300.0 && dist <= 500.0:
			return 3600, 5000
		default:
			return 2700, 3750
		}
	}

	fdist = func(dist float64) float64 {
		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 50
		case dist > 5.0 && dist <= 100.0:
			return 100
		case dist > 100.0 && dist <= 200.0:
			return 200
		case dist > 200.0 && dist <= 300.0:
			return 300
		case dist > 300.0 && dist <= 500.0:
			return 500
		default:
			return -1
		}
	}
	max_pts = 20
	office_time = 15 * 60
}

func main() {
	go Start()
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Start() {
	wh := Point{
		Lat: 55.58783,
		Lng: 37.226277,
	}

	ptsMap := map[Point]int{}
	pts := []Point{wh}

	f, err := os.Open("vnukovo.csv")
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
		office_id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Fatal(err)
		}
		avg_shk, err := strconv.Atoi(row[4])
		if err != nil {
			log.Fatal(err)
		}
		lat, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		lng, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		ll := Point{
			Id:  office_id,
			Shk: avg_shk,
			Lat: lat,
			Lng: lng,
		}
		pts = append(pts, ll)

		ptsMap[ll] += 1
	}

	// матрицы расстояний
	l := len(pts)
	dm := make([][]float64, l)
	tm := make([][]float64, l)

	for i := 0; i < l; i++ {
		dm[i] = make([]float64, l)
		tm[i] = make([]float64, l)
	}

	max_len := 500
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
			ttm, tdm := GetMatrices2(pts, src, dst) // times and distances
			for iii, row := range tdm {
				for jjj, p := range row {
					dm[iii+j*max_len][jjj+jj*max_len] = p
					tm[iii+j*max_len][jjj+jj*max_len] = ttm[iii][jjj]
				}
			}
		}
	}

	fmt.Println(len(dm))

	initialRoutes := initialGreedy(pts, dm)

	currSolution := make([][]int, len(initialRoutes))
	for i, r := range initialRoutes {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}
	var (
		newFitness  float64
		oldFitness  float64
		oldOvershk  int
		oldUndershk int
		newOvershk  int
		newUndershk int
	)
	oldFitness, oldOvershk, oldUndershk = fitness(initialRoutes, pts, dm)
	fmt.Println("initial status:", len(currSolution), oldFitness, oldOvershk)
	for i := 0; i < 1; i++ {
		fmt.Println("iteration:", i)

		currSolution = sa(currSolution, pts, dm, 0.99)
		newFitness, newOvershk, newUndershk = fitness(currSolution, pts, dm)
		oldFitness = newFitness
		oldOvershk = newOvershk
		oldUndershk = newUndershk

		fmt.Println("main cycle update:", len(currSolution), oldFitness, oldOvershk, oldUndershk)
	}

	currSolution = sa(currSolution, pts, dm, 0.9999)
	newFitness, newOvershk, newUndershk = fitness(currSolution, pts, dm)

	prepareGeoJson(currSolution, pts, dm, tm)
	prepareFile(currSolution, pts, tm)

	fmt.Println(len(currSolution), newFitness, newOvershk, newUndershk)
}
