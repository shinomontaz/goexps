package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func prepareFile(routes [][]int, pts []Point, tm [][]float64) {
	file, err := os.Create("result.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	// this defines the header value and data values for the new csv file

	items := make([][]string, 0, len(routes)+1)
	items = append(items, []string{"office", "route", "time", "lat", "lng"})
	for rid, r := range routes {
		dt := 0.0
		for i := 1; i < len(r); i++ {
			dt += tm[r[i-1]][r[i]]
			items = append(items, []string{
				fmt.Sprintf("%d", pts[r[i]].Id),
				fmt.Sprintf("%d", rid+1),
				fmt.Sprintf("%f", dt),
				fmt.Sprintf("%f", pts[r[i]].Lat),
				fmt.Sprintf("%f", pts[r[i]].Lng),
			})
			dt += office_time
		}
	}
	if writer.WriteAll(items) != nil {
		panic("!")
	}
}

type DataRouted struct {
	Pts      []RoutePoint `json:"pts"`
	Routes   []RouteJson  `json:"routes"`
	Cost     float64      `json:"cost"`
	Time     float64      `json:"time"`
	Distance float64      `json:"distance"`
	Wh       int          `json:"wh"`
	BoxTime  float64      `json:"box_time"`
}

type RoutePoint struct {
	OfficeId int     `json:"office_id"`
	Route    int     `json:"route"`
	Time     string  `json:"time"`
	Shk      int     `json:"shk"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type RouteJson struct {
	Points   []int   `json:"points"`
	Time     float64 `json:"time"`
	Shk      int     `json:"shk"`
	ShkLim   string  `json:"shk_lim"`
	Distance float64 `json:"distance"`
	Geojson  Geojson `json:"geojson"`
}

func prepareGeoJson(routes [][]int, pts []Point, dm, tm [][]float64) {
	res := DataRouted{
		Routes: make([]RouteJson, 0, len(routes)),
		Pts:    make([]RoutePoint, 0, len(pts)),
	}

	totalDist, _, _ := fitness(routes, pts, dm) // cost, time, overbox, overshift float64
	res.Cost = totalDist
	for rId, r := range routes {
		if len(r) <= 1 {
			continue
		}
		route := routeGeoJson(r, pts, dm, tm)
		res.Routes = append(res.Routes, route)
		res.Distance += route.Distance
		res.Time += route.Time

		for i := 1; i < len(r); i++ {
			res.Pts = append(res.Pts, RoutePoint{
				OfficeId: pts[r[i]].Id,
				Shk:      pts[r[i]].Shk,
				Lat:      pts[r[i]].Lat,
				Lng:      pts[r[i]].Lng,
				Route:    rId,
			})
		}
	}

	file, _ := json.MarshalIndent(res, "", " ")

	_ = os.WriteFile("result.json", file, 0644)
}

// DataRouted.Pts will be changed here: we add route id, num and time
func routeGeoJson(route []int, pts []Point, dm, tm [][]float64) RouteJson {
	res := RouteJson{
		Points: route,
	}

	routepts := []Point{pts[0]}
	dist := 0.0
	time := 0.0

	for i := 1; i < len(route); i++ {
		routepts = append(routepts, pts[route[i]])
		dist += dm[route[i-1]][route[i]]
		time += tm[route[i-1]][route[i]] + office_time
		res.Shk += pts[route[i]].Shk

	}

	if len(route) > 1 {
		// добаляем время на возврат на склад
		dist += dm[route[len(route)-1]][0]
		time += tm[route[len(route)-1]][0]
	}

	_, _, gj, _ := GetRoute(routepts)

	res.Geojson = gj
	shk_min, shk_max := fshk(dist / 1000.0)
	res.ShkLim = fmt.Sprintf("%d/%d", shk_min, shk_max)
	res.Time = time
	res.Distance = dist

	return res
}
