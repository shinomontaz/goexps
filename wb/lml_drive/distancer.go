package main

import (
	"encoding/json"
	"fmt"
	"lml-drive/types"
	"log"
	"net/http"
	"strings"
)

type Geojson struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type OSRMApiTableResp struct {
	Code         string      `json:"code"`
	Distances    [][]float64 `json:"distances"`
	Durations    [][]float64 `json:"durations"`
	Destinations []struct {
		Hint     string    `json:"hint"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"destinations"`
	Sources []struct {
		Hint     string    `json:"hint"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"sources"`
}

type OSRMApiRouteResp struct {
	Code   string `json:"code"`
	Routes []struct {
		Geometry Geojson `json:"geometry"`
		Legs     []struct {
			Steps    []interface{} `json:"steps"`
			Distance float64       `json:"distance"`
			Duration float64       `json:"duration"`
			Summary  string        `json:"summary"`
			Weight   float64       `json:"weight"`
		} `json:"legs"`
		Distance   float64 `json:"distance"`
		Duration   float64 `json:"duration"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
	} `json:"routes"`
	Waypoints []struct {
		Hint     string    `json:"hint"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"waypoints"`
}

// GetRoute returns distance, duration, geojson and error
func GetRoute(points []types.Point) (float64, float64, Geojson, error) {
	qsParts := make([]string, 0)
	for _, p := range points {
		qsParts = append(qsParts, fmt.Sprintf("%f,%f", p.Lng, p.Lat))
	}

	resp, err := http.Get(fmt.Sprintf("%s/route/v1/driving/%s?geometries=geojson&overview=full", OsrmURL, strings.Join(qsParts, ";")))
	if err != nil {
		panic(err)
	}
	apiResp := OSRMApiRouteResp{}
	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, 0, Geojson{}, fmt.Errorf("can't unmarshal osrm api resp: %s", err)
	}
	if len(apiResp.Routes) == 0 {
		return 0, 0, Geojson{}, fmt.Errorf("not route")
	}

	return apiResp.Routes[0].Distance, apiResp.Routes[0].Duration, apiResp.Routes[0].Geometry, nil
}

func GetMatricesChunked(pts []types.Point) ([][]float64, [][]float64) {
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
	return dm, tm
}

func GetMatrices(pts []types.Point) ([][]float64, [][]float64) {
	qsParts := make([]string, 0)
	for _, p := range pts {
		qsParts = append(qsParts, fmt.Sprintf("%f,%f", p.Lng, p.Lat))
	}

	qs := fmt.Sprintf("%s/table/v1/driving/%s?annotations=duration,distance", OsrmURL, strings.Join(qsParts, ";"))

	resp, err := http.Get(qs)

	if err != nil {
		log.Fatalf("can't call osrm api: %s", err)
	}
	apiResp := OSRMApiTableResp{}
	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Fatalf("can't unmarshal osrm api resp: %s", err)
	}

	return apiResp.Durations, apiResp.Distances
}

// return times and distances
func GetMatrices2(pts []types.Point, src, dst []int) ([][]float64, [][]float64) {
	qsParts := make([]string, 0)
	for _, p := range pts {
		qsParts = append(qsParts, fmt.Sprintf("%f,%f", p.Lng, p.Lat))
	}

	srcs := strings.Trim(strings.Replace(fmt.Sprint(src), " ", ";", -1), "[]")
	dsts := strings.Trim(strings.Replace(fmt.Sprint(dst), " ", ";", -1), "[]")

	qs := fmt.Sprintf("%s/table/v1/driving/%s?annotations=duration,distance&sources=%s&destinations=%s", OsrmURL, strings.Join(qsParts, ";"), srcs, dsts)

	resp, err := http.Get(qs)

	if err != nil {
		log.Fatalf("can't call osrm api: %s", err)
	}
	apiResp := OSRMApiTableResp{}
	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Fatalf("can't unmarshal osrm api resp: %s", err)
	}

	return apiResp.Durations, apiResp.Distances
}
