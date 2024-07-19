package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// returns distance and duration
func GetRoute(p, q LatLong) (float64, float64, error) {
	if p == q {
		return 0, 0, nil
	}

	qsParts := fmt.Sprintf("%f,%f;%f,%f", p.Lng, p.Lat, q.Lng, q.Lat)
	url := fmt.Sprintf("%s/route/v1/driving/%s?overview=false", OsrmURL, qsParts)
	resp, err := http.Get(url)
	apiResp := OSRMApiRouteResp{}
	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, 0, fmt.Errorf("can't unmarshal osrm api resp: %s", err)
	}
	if len(apiResp.Routes) == 0 {
		return 0, 0, fmt.Errorf("not route")
	}

	return apiResp.Routes[0].Distance, apiResp.Routes[0].Duration, nil
}

func GetMatrices(pts []LatLong) ([][]float64, [][]float64) {
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
func GetMatrices2(pts []LatLong, src, dst []int) ([][]float64, [][]float64) {
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
