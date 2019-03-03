package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func prepare(filein, fileout string) {

	hubs := make(map[string]*Order)
	destinations := make(map[string]*Order)

	f, err := os.Open(filein)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	tasks := make(map[string]*Task)

	failing_addresses := make([]string, 0, 10)

	// 0			1				2	3	4			5
	// hub_name		hub_address		-	-	dest_name	dest_address
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if _, exists := hubs[line[0]]; !exists {
			hubs[line[0]] = &Order{
				Address: line[1],
				Guid:    line[0],
			}

			time.Sleep(1 * time.Second)
			x, y, err := Geocode(hubs[line[0]].Address, true)
			if err != nil {
				panic(err)
			}

			coords := &LatLong{
				Lat:  x,
				Long: y,
			}

			if coords == nil {
				failing_addresses = append(failing_addresses, fmt.Sprintf("%s - %s", hubs[line[0]].Guid, hubs[line[0]].Address))
				continue
			}

			hubs[line[0]].Coords = coords

			tasks[line[0]] = &Task{
				Hub:    hubs[line[0]],
				Points: make([]*Order, 0),
			}
		}

		if _, exists := destinations[line[4]]; !exists {
			destinations[line[4]] = &Order{
				Address: line[5],
				Guid:    line[4],
			}

			time.Sleep(1 * time.Second)
			x, y, err := Geocode(destinations[line[4]].Address, true)
			if err != nil {
				panic(err)
			}

			coords := &LatLong{
				Lat:  x,
				Long: y,
			}

			if coords == nil {
				failing_addresses = append(failing_addresses, fmt.Sprintf("%s - %s", destinations[line[4]].Guid, destinations[line[4]].Address))
				continue
			}

			destinations[line[4]].Coords = coords

			if _, exists := tasks[line[0]]; exists {
				tasks[line[0]].Points = append(tasks[line[0]].Points, destinations[line[4]])
			}
		}

	}

	if len(failing_addresses) > 0 {
		for _, entry := range failing_addresses {
			fmt.Println(entry)
		}
	}

	file, err := os.Create(fileout)
	if err != nil {
		file, err = os.Open(fileout)
		if err != nil {
			panic(err)
		}
	}

	defer file.Close()

	for title, task := range tasks {
		fmt.Println(title)
		for _, point := range task.Points {
			fmt.Println("point:", point)
		}
	}

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for hubName, task := range tasks {
		for _, point := range task.Points {
			value := []string{
				hubName,
				task.Hub.Address,
				fmt.Sprintf("%f", task.Hub.Coords.Lat),
				fmt.Sprintf("%f", task.Hub.Coords.Long),
				point.Guid,
				point.Address,
				fmt.Sprintf("%f", point.Coords.Lat),
				fmt.Sprintf("%f", point.Coords.Long),
			}
			err := writer.Write(value)
			if err != nil {
				panic("!")
			}
		}
	}
}
