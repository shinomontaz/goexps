package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/shinomontaz/goexps/Taboo/types"
)

var orders = make([]*types.Order, 0)

func init() {
	file, _ := os.Open("./data/orders.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&orders)
	if err != nil {
		log.Fatal(err)
	}
}
