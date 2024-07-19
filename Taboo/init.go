package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

var ords = make([]LatLng, 0)
var rs *rand.Rand
var dm [][]float64

func init() {
	file, _ := os.Open("./data/orders.json")
	decoder := json.NewDecoder(file)
	orders := []*Order{}
	err := decoder.Decode(&orders)
	if err != nil {
		log.Fatal(err)
	}

	ords = make([]LatLng, len(orders))
	for i, o := range orders {
		ords[i] = LatLng{
			Lat: o.Lat,
			Lng: o.Long,
		}
	}

	rs = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	dm = make([][]float64, len(ords))

	for i := 0; i < len(ords); i++ {
		dm[i] = make([]float64, len(ords))
		for j := 0; j < len(ords); j++ {
			dm[i][j] = ords[i].Distance(ords[j])
		}
	}

}
