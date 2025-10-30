package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

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

	pts := []Point{}
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
	mst, cost := buildMst(pts, dm)
	fmt.Println(mst, cost)
}
