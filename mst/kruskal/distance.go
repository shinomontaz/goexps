package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type osrmApiTableResp struct {
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

const osrmUrl = "http://osrm-gazelle.wbdispatch.k8s.prod-xc"

func getMatrix(pts []Point) [][]float64 {
	qsParts := make([]string, 0)
	for _, p := range pts {
		qsParts = append(qsParts, fmt.Sprintf("%f,%f", p.Coord.Lon, p.Coord.Lat))
	}

	qs := fmt.Sprintf("%s/table/v1/driving/%s", osrmUrl, strings.Join(qsParts, ";"))

	resp, err := http.Get(qs)

	if err != nil {
		fmt.Printf("can't call osrm api: %s", err)
		panic("!")
	}
	apiResp := osrmApiTableResp{}
	if json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("can't unmarshal osrm api resp: %s", err)
		panic("!")
	}

	return apiResp.Durations
}
