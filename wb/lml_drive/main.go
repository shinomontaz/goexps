package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"lml-drive/rand"
	"lml-drive/types"
	"os"
	"strconv"

	"log"

	"lml-drive/mst"
	"net/http"
	_ "net/http/pprof"
)

type dist_shk_f func(dist float64) (int, int)
type dist_max_f func(dist float64) float64

var (
	OsrmURL string
	fshk    dist_shk_f
	fdist   dist_max_f
	fmnshk  func(d float64) int
	fmxshk  func(d float64) int
	fbushk  func(d float64) int

	max_pts     int
	office_time float64
	result_name string
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
	OsrmURL = "http://osrm-gazelle.wbdispatch.k8s.prod-xc" // http://91.228.153.227:5000" //
	rand.Init(1 << 24)
	fshk = func(dist float64) (int, int) {
		dist /= 1000.0

		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 4500, 6250
		case dist > 50.0 && dist <= 100.0:
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
		dist /= 1000.0

		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 50
		case dist > 50.0 && dist <= 100.0:
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

	fmnshk = func(dist float64) int {
		dist /= 1000.0
		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 4500
		case dist > 50.0 && dist <= 100.0:
			return 4500
		case dist > 100.0 && dist <= 200.0:
			return 4320
		case dist > 200.0 && dist <= 300.0:
			return 3780
		case dist > 300.0 && dist <= 500.0:
			return 3600
		default:
			return 2700
		}
	}
	fmxshk = func(dist float64) int {
		dist /= 1000.0

		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 6250
		case dist > 50.0 && dist <= 100.0:
			return 6250
		case dist > 100.0 && dist <= 200.0:
			return 6000
		case dist > 200.0 && dist <= 300.0:
			return 5250
		case dist > 300.0 && dist <= 500.0:
			return 5000
		default:
			return 3750
		}
	}
	fbushk = func(dist float64) int {
		dist /= 1000.0
		switch {
		case dist >= 0.0 && dist <= 50.0:
			return 1
		case dist > 50.0 && dist <= 100.0:
			return 2
		case dist > 100.0 && dist <= 200.0:
			return 3
		case dist > 200.0 && dist <= 300.0:
			return 4
		case dist > 300.0 && dist <= 500.0:
			return 5
		default:
			return 6
		}
	}
	max_pts = 20
	office_time = 15 * 60
	mst.Init(fmnshk, fmxshk, fbushk)
}

func main() {
	//	result_name = "result_irkutsk"
	result_name = "result_krylovskaya"
	go Start()
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Start() {
	// wh := Point{
	// 	Lat: 55.58783,
	// 	Lng: 37.226277,
	// }

	// wh := Point{ // kaluga
	// 	Lat: 54.483186,
	// 	Lng: 36.218078,
	// }

	// wh := Point{ // rostov 151 & 152
	// 	Lat: 47.367953,
	// 	Lng: 39.71672,
	// }

	wh := types.Point{ // Крыловская
		Lat: 46.312237,
		Lng: 39.946724,
	}
	// wh := types.Point{ // Irkutsk
	// 	Lat: 52.3196490,
	// 	Lng: 104.2482280,
	// }

	pts := []types.Point{wh}

	//	f, err := os.Open("routes_for_krasnodar_Igor.csv")
	//	f, err := os.Open("Krasnodar.csv")
	f, err := os.Open("Krylovskaya.csv")
	//	f, err := os.Open("Irkutsk.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	//	csvReader.Comma = ';'

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
		avg_shk, err := strconv.Atoi(row[3])
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

		ll := types.Point{
			Id:  office_id,
			Shk: avg_shk,
			Lat: lat,
			Lng: lng,
		}
		pts = append(pts, ll)
	}

	var table [][]string

	resGeoJson := DataRouted{
		Routes: make([]RouteJson, 0),
		Pts:    make([]RoutePoint, 0, len(pts)+1),
	}

	resGeoJson.Pts = append(resGeoJson.Pts, RoutePoint{
		Lat: wh.Lat,
		Lng: wh.Lng,
	})

	var (
		total_ovshk int
		total_unshk int
	)

	dm, tm := GetMatricesChunked(pts)
	//	clusters := clusterize(start_num_clusters, pts)
	clusters := mst.Do(pts, dm, tm)
	sum := 0
	for _, c := range clusters {
		fmt.Println(len(c))
		sum += len(c)
	}
	fmt.Println("total: ", sum)

	t := make([][]string, len(clusters))
	t = append(t, []string{"office", "route", "time", "lat", "lng"})
	last_route_id := 0
	for _, c := range clusters {
		fmt.Println("cluster size: ", len(c))
		gj, t, ovshk, unshk := solve(c, last_route_id)
		table = append(table, t...)
		last_route_id += len(gj.Routes)

		offset := len(resGeoJson.Pts) - 1
		for i := 0; i < len(gj.Pts); i++ {
			if len(resGeoJson.Pts) > 0 && i == 0 {
				continue
			}
			gj.Pts[i].Route += len(resGeoJson.Routes)
			resGeoJson.Pts = append(resGeoJson.Pts, gj.Pts[i])
		}
		for _, r := range gj.Routes {
			for i := 1; i < len(r.Points); i++ {
				r.Points[i] += offset
			}
			resGeoJson.Routes = append(resGeoJson.Routes, r)
		}

		resGeoJson.Cost += gj.Cost
		resGeoJson.Distance += gj.Distance
		total_ovshk += ovshk
		total_unshk += unshk
	}

	fmt.Println("resGeoJson.Pts: ", len(resGeoJson.Pts), len(resGeoJson.Routes))

	file, _ := json.MarshalIndent(resGeoJson, "", " ")
	_ = os.WriteFile(result_name+".json", file, 0644)

	ff, err := os.Create(result_name + ".csv")
	if err != nil {
		panic(err)
	}
	defer ff.Close()

	writer := csv.NewWriter(ff)
	defer writer.Flush()
	if writer.WriteAll(table) != nil {
		panic("!")
	}

}
