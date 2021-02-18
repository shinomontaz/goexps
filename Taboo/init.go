package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/shinomontaz/goexps/types"
)

var orders = make([]*types.Order, 0)

func init() {
	rand.Seed(time.Now().UnixNano())
	file, _ := os.Open("./data/orders.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&orders)
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range orders {
		o.Coord = &types.LatLng{
			Lat: o.Lat,
			Lng: o.Long,
		}
	}
}
