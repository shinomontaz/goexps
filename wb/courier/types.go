package main

type LatLong struct {
	Lat float64
	Lng float64
}

type Point struct {
	Uuid   string
	Coord  LatLong
	Vol    float64
	Weight float64
	Cid    int
}

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
