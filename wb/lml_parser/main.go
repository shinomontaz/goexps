package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

var filename = "tares_12032024.csv"

type Order struct {
	ID      int     `json:"id"`
	Zone    int     `json:"zone"`
	Route   int     `json:"route"`
	Lng     float64 `json:"lng"`
	Lat     float64 `json:"lat"`
	WhLng   float64 `json:"whLng"`
	WhLat   float64 `json:"whLat"`
	Task_id int     `json: "task"`
}

func main() {
	file, err := os.Open(filename)
	defer file.Close()

	reader := csv.NewReader(file)
	tares, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("cannot decode orders mp json: %v", err)
	}
	whlat := 55.386140000
	whlng := 37.583393000

	var (
		task_id int
		tare_id int
		// src_office_id int
		// dst_office_id int
		route_id   int
		car_number string
		lng        float64
		lat        float64
	)
	cars := make(map[string]int)
	car_orders := make(map[string][]Order)
	car_tasks := make(map[string]map[int]struct{})

	// группировка: task_id ( для одного task_id должна быть 1 машина и 1 маршрут )
	// если task_id разный, но src_office_id и номер машины одинаковый - то объединяем
	for _, row := range tares[1:] {
		//		fmt.Println(i + 1)
		//		"app","report_d","task_id","tare_id","src_office_id","dst_office_id","route_id","car_number","delivery_time_hours","office_name","city_name","full_address","lat","long"
		if task_id, err = strconv.Atoi(row[2]); err != nil {
			log.Fatalf("cannot convert task_id to int: %v", err)
			return
		}
		if tare_id, err = strconv.Atoi(row[3]); err != nil {
			log.Fatalf("cannot convert tare_id to int: %v", err)
			return
		}
		if row[0] == "WBDRIVE" {
			route_id = -1
		} else {
			if route_id, err = strconv.Atoi(row[6]); err != nil {
				log.Fatalf("cannot convert route_id to int: %v", err)
				return
			}
		}
		car_number = row[7]

		if _, ok := cars[car_number]; !ok {
			cars[car_number] = len(cars)
			car_orders[car_number] = []Order{}
			car_tasks[car_number] = make(map[int]struct{})
		}
		if lat, err = strconv.ParseFloat(row[12], 64); err != nil {
			log.Fatalf("cannot convert lat to float64: %v", err)
			return
		}
		if lng, err = strconv.ParseFloat(row[13], 64); err != nil {
			log.Fatalf("cannot convert lng to float64: %v", err)
			return
		}

		if _, ok := car_tasks[car_number][task_id]; !ok {
			car_tasks[car_number][task_id] = struct{}{}
		}

		car_orders[car_number] = append(car_orders[car_number], Order{
			ID:      tare_id,
			Lat:     lat,
			Lng:     lng,
			Zone:    route_id,
			Route:   cars[car_number],
			WhLat:   whlat,
			WhLng:   whlng,
			Task_id: task_id,
		})
	}

	fmt.Println(len(car_orders))

	for cn, ct := range car_tasks {
		if len(ct) > 1 {
			fmt.Println(cn, ct)
		}
	}
}
