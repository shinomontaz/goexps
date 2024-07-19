package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type OsrmRouteResp struct {
	Code   string `json:"code"`
	Routes []struct {
		Geometry string `json:"geometry"`
		Legs     []struct {
			Steps []struct {
				Geometry string `json:"geometry"`
				Maneuver struct {
					BearingAfter  int       `json:"bearing_after"`
					BearingBefore int       `json:"bearing_before"`
					Location      []float64 `json:"location"`
					Modifier      string    `json:"modifier"`
					Type          string    `json:"type"`
				} `json:"maneuver"`
				Mode          string `json:"mode"`
				DrivingSide   string `json:"driving_side"`
				Name          string `json:"name"`
				Intersections []struct {
					Out      int       `json:"out"`
					Entry    []bool    `json:"entry"`
					Bearings []int     `json:"bearings"`
					Location []float64 `json:"location"`
					In       int       `json:"in,omitempty"`
				} `json:"intersections"`
				Weight       float64 `json:"weight"`
				Duration     float64 `json:"duration"`
				Distance     float64 `json:"distance"`
				Ref          string  `json:"ref,omitempty"`
				Destinations string  `json:"destinations,omitempty"`
			} `json:"steps"`
			Summary  string  `json:"summary"`
			Weight   float64 `json:"weight"`
			Duration float64 `json:"duration"`
			Distance float64 `json:"distance"`
		} `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	} `json:"routes"`
	Waypoints []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"waypoints"`
}

func main() {
	dat, err := os.ReadFile("data.txt")
	if err != nil {
		panic(err)
	}

	var route OsrmRouteResp
	err = json.Unmarshal([]byte(dat), &route)
	if err != nil {
		panic(err)
	}

	// print all
	for _, l := range route.Routes[0].Legs {
		// print speed
		fmt.Printf("%s: %f km/h\n", l.Summary, (l.Distance/1000.0)/(l.Duration/3600.0))
		for _, s := range l.Steps {
			// print speed
			fmt.Printf("\t%s: %f km/h\n", s.Destinations, (s.Distance/1000.0)/(s.Duration/3600.0))
		}

	}
}
