package main

import (
	"courier/rand"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"log"

	"net/http"
	_ "net/http/pprof"
)

var (
	OsrmURL string

	max_pts     int
	office_time float64
	rc          routeChecker
)

func init() {
	OsrmURL = "http://osrm-gazelle.wbdispatch.k8s.prod-xc"
	rand.Init(1 << 24)

	max_pts = 10
	office_time = 5 * 60
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
		Coord: LatLong{
			Lat: 43.235406,
			Lng: 76.892700,
		},
		Cid: 0,
	}

	pts := []Point{wh}

	f, err := os.Open("data.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '|'

	_, err = csvReader.Read() // read header
	if err != nil {
		log.Fatal(err)
	}

	coordMap := map[LatLong]int{wh.Coord: 0}

	i := 1
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lat, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		lng, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			log.Fatal(err)
		}

		crd := LatLong{
			Lat: lat,
			Lng: lng,
		}
		if _, ok := coordMap[crd]; !ok {
			coordMap[crd] = len(coordMap)
		}

		ll := Point{
			Uuid:  fmt.Sprintf("%d", i),
			Coord: crd,
			Cid:   coordMap[crd],
		}
		pts = append(pts, ll)
		i++
	}

	upts := make([]LatLong, len(coordMap))
	for k, v := range coordMap {
		upts[v] = k
	}

	// матрицы расстояний
	l := len(upts)
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
			ttm, tdm := GetMatrices2(upts, src, dst) // times and distances
			for iii, row := range tdm {
				for jjj, p := range row {
					dm[iii+j*max_len][jjj+jj*max_len] = p
					tm[iii+j*max_len][jjj+jj*max_len] = ttm[iii][jjj]
				}
			}
		}
	}

	rc = routeChecker{
		LimitDist:  30 * 1000,
		LimitNum:   11,
		seenCoords: make(map[LatLong]struct{}),
	}

	initialRoutes := initialGreedy(pts, dm)

	currSolution := make([][]int, len(initialRoutes))
	for i, r := range initialRoutes {
		currSolution[i] = make([]int, len(r))
		copy(currSolution[i], r)
	}

	currSolution = sa(currSolution, pts, dm, tm, 0.9999)

	prepareGeoJson(currSolution, pts, dm, tm)

	fmt.Println(len(currSolution))
}
