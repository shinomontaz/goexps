package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func saveMatrix(dm [][]float64, tm [][]float64) {
	file, err := os.Create("matrix.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	// this defines the header value and data values for the new csv file

	//	str_row := make([]string, 0, len(dm))
	items := make([][]string, len(dm))
	for i, row := range dm {
		items[i] = make([]string, len(row))
		for j, d := range row {
			items[i][j] = fmt.Sprintf("(%f;%f)", d, tm[i][j])
		}
	}
	if writer.WriteAll(items) != nil {
		panic("!")
	}
}

type DataRouted struct {
	Pts        []RoutePoint `json:"pts"`
	Routes     []RouteJson  `json:"routes"`
	Cost       float64      `json:"cost"`
	Time       float64      `json:"time"`
	Distance   float64      `json:"distance"`
	MeanDist   float64      `json:"mean_distance"`
	MeanTime   float64      `json:"mean_time"`
	MeanPoints float64      `json:"mean_points"`
	Wh         int          `json:"wh"`
	BoxTime    float64      `json:"box_time"`
}

type RoutePoint struct {
	Uuid  string  `json:"office_id"`
	Route int     `json:"route"`
	Time  string  `json:"time"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}

type RouteJson struct {
	Points   []int   `json:"points"`
	Time     float64 `json:"time"`
	Distance float64 `json:"distance"`
	Geojson  Geojson `json:"geojson"`
}

func prepareGeoJson(routes [][]int, pts []Point, dm, tm [][]float64) {
	res := DataRouted{
		Routes: make([]RouteJson, 0, len(routes)),
		Pts:    make([]RoutePoint, 0, len(pts)),
	}

	res.Pts = append(res.Pts, RoutePoint{
		Uuid: pts[0].Uuid,
		Lat:  pts[0].Coord.Lat,
		Lng:  pts[0].Coord.Lng,
	})

	totalDist := fitness(routes, pts, dm) // cost, time, overbox, overshift float64
	res.Cost = totalDist
	for rId, r := range routes {
		if len(r) <= 1 {
			continue
		}
		route := routeGeoJson(r, pts, dm, tm)
		res.Routes = append(res.Routes, route)
		res.Distance += route.Distance
		res.Time += route.Time

		for _, p := range r {
			res.Pts = append(res.Pts, RoutePoint{
				Uuid:  pts[p].Uuid,
				Lat:   pts[p].Coord.Lat,
				Lng:   pts[p].Coord.Lng,
				Route: rId,
			})
		}
	}
	/*
		MeanPoints float64      `json:"mean_points"`
	*/
	res.MeanDist = res.Distance / float64(len(res.Routes))
	res.MeanTime = res.Time / float64(len(res.Routes))
	res.MeanPoints = float64(len(pts)-1) / float64(len(res.Routes))

	file, _ := json.MarshalIndent(res, "", " ")

	_ = os.WriteFile("result_courier.json", file, 0644)
}

// DataRouted.Pts will be changed here: we add route id, num and time
func routeGeoJson(route []int, pts []Point, dm, tm [][]float64) RouteJson {
	res := RouteJson{
		Points: route,
	}

	routepts := []LatLong{pts[0].Coord}
	dist := 0.0
	time := 0.0

	seenCoord := map[LatLong]struct{}{}
	for i := 1; i < len(route); i++ {
		if _, ok := seenCoord[pts[route[i]].Coord]; !ok {
			routepts = append(routepts, pts[route[i]].Coord)
		}
		dist += dm[pts[route[i-1]].Cid][pts[route[i]].Cid]
		time += tm[pts[route[i-1]].Cid][pts[route[i]].Cid]
		seenCoord[pts[route[i]].Coord] = struct{}{}
	}

	if len(route) > 1 {
		// добаляем время на возврат на склад
		dist += dm[pts[route[len(route)-1]].Cid][0]
		time += tm[pts[route[len(route)-1]].Cid][0]
	}

	_, _, gj, _ := GetRoute(routepts)

	res.Geojson = gj
	res.Time = time
	res.Distance = dist

	return res
}
