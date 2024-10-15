package main

import (
	"encoding/csv"
	"fmt"
	"lml-drive/types"
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

func prepareFile(routes [][]int, pts []types.Point, tm [][]float64, last_route_id int) [][]string {
	// file, err := os.Create("result_k.csv")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// writer := csv.NewWriter(file)
	// defer writer.Flush()
	// this defines the header value and data values for the new csv file

	items := make([][]string, 0, len(routes))
	for rid, r := range routes {
		dt := 0.0
		for i := 1; i < len(r); i++ {
			dt += tm[r[i-1]][r[i]]
			items = append(items, []string{
				fmt.Sprintf("%d", pts[r[i]].Id),
				fmt.Sprintf("%d", rid+1+last_route_id),
				fmt.Sprintf("%f", dt),
				fmt.Sprintf("%f", pts[r[i]].Lat),
				fmt.Sprintf("%f", pts[r[i]].Lng),
			})
			dt += office_time
		}
	}
	return items
	// if writer.WriteAll(items) != nil {
	// 	panic("!")
	// }
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

func prepareGeoJson(routes [][]int, pts []types.Point, dm, tm [][]float64) DataRouted {
	res := DataRouted{
		Routes: make([]RouteJson, 0, len(routes)),
		Pts:    make([]RoutePoint, 0, len(pts)),
	}

	res.Pts = append(res.Pts, RoutePoint{
		OfficeId: pts[0].Id,
		Shk:      0,
		Lat:      pts[0].Lat,
		Lng:      pts[0].Lng,
	})

	totalDist, _, _, _ := fitness(routes, pts, dm) // cost, time, overbox, overshift float64
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
	/*
		MeanPoints float64      `json:"mean_points"`
	*/
	res.MeanDist = res.Distance / float64(len(res.Routes))
	res.MeanTime = res.Time / float64(len(res.Routes))
	res.MeanPoints = float64(len(pts)-1) / float64(len(res.Routes))

	return res
	// file, _ := json.MarshalIndent(res, "", " ")

	// _ = os.WriteFile("result_k.json", file, 0644)
}

// DataRouted.Pts will be changed here: we add route id, num and time
func routeGeoJson(route []int, pts []types.Point, dm, tm [][]float64) RouteJson {
	res := RouteJson{
		Points: route,
	}

	routepts := []types.Point{pts[0]}
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
	shk_min, shk_max := fshk(dist)
	res.ShkLim = fmt.Sprintf("%d/%d", shk_min, shk_max)
	res.Time = time
	res.Distance = dist

	return res
}
