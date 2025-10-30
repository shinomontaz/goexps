package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

var clusters int
var cluster_num [2]int

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed reading CSV: %s", err)
	}

	pts := []Point{
		{
			Id: 0,
			Coord: LatLon{
				Lat: 55.637672,
				Lon: 37.697581,
			},
		},
	}
	//lat;long;office_id;shk_cnt
	i := 1
	for _, record := range records[1:] { // Skip the header
		lat, _ := strconv.ParseFloat(record[0], 64)
		lon, _ := strconv.ParseFloat(record[1], 64)
		officeID, _ := strconv.Atoi(record[2])
		shk, _ := strconv.Atoi(record[3])

		pts = append(pts, Point{
			Id: int64(officeID),
			Coord: LatLon{
				Lat: lat,
				Lon: lon,
			},
			Shk: shk,
		})
		i++
	}

	dm := getMatrix(pts)

	clusters = 12
	meanPointsNum := len(pts) / clusters
	cluster_num = [2]int{meanPointsNum - int(0.1*float64(meanPointsNum)), meanPointsNum + int(0.1*float64(meanPointsNum))}

	clusters, _ := buildMst(pts[1:], dm)
	fmt.Println(clusters)
	initial := make([][]int, len(clusters))
	for i, c := range clusters {
		initial[i] = []int{0}
		for _, p := range c {
			initial[i] = append(initial[i], p+1)
		}
	}
	fmt.Println(initial)

	solutions, cost := sa(initial, dm)

	fmt.Println(solutions, cost)

	var targetData TargetData

	// Convert Routes
	for i, route := range solutions {
		var targetRoute TargetRoute

		// Convert Points
		for j, p := range route {
			var targetPoint TargetPoint
			targetPoint.SupplierID = 0
			targetPoint.OrderID = pts[p].Id
			targetPoint.Route = j
			targetPoint.Num = i
			targetPoint.Lat = pts[p].Coord.Lat
			targetPoint.Lon = pts[p].Coord.Lon

			targetRoute.Points = append(targetRoute.Points, targetPoint.OrderID)
			targetPoint.PointId = pts[p].Id
			targetData.Pts = append(targetData.Pts, targetPoint)
		}

		targetData.Routes = append(targetData.Routes, targetRoute)
	}

	outputData, err := json.MarshalIndent(targetData, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)

		return
	}

	if os.WriteFile("pnd_sa_output.json", outputData, 0644) != nil { //nolint:all
		fmt.Println("Error writing output file:", err)
	}
	// to mapjson!
}
