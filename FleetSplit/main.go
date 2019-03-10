package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	tags := []*Tag{
		&Tag{
			Id:   1,
			Name: "Ювелирка",
		},
		&Tag{
			Id:   2,
			Name: "КГТ",
		},
	}

	zonesNum := 30
	zones := make([]*Zone, 0, zonesNum)
	orders := make([]*Order, 0, 5000)
	for i := 0; i < 5000; i++ {
		orders = append(orders, &Order{
			Zone: int64(rand.Intn(5000) % zonesNum),
			Parameters: &Limits{
				Weight: float64(rand.Intn(5) + (rand.Intn(10)%5)*5),
			},
		})
	}

	ordersByZones := make(map[int][]*Order, 0)
	for _, ord := range orders {
		ordersByZones[int(ord.Zone)] = append(ordersByZones[int(ord.Zone)], ord)
	}

	for i := 0; i < zonesNum; i++ {
		zoneWeight := 0.0
		for _, ord := range ordersByZones[i] {
			zoneWeight += ord.Parameters.Weight
		}

		zones = append(zones, &Zone{
			ZoneId: i,
			Weight: zoneWeight,
		})
	}

	fleet := generateFleet(100, 50, zones, tags)

	hardfleet := make(map[int][]*Courier)
	couriers := make([]*Courier, 0, len(fleet))

	for _, courier := range fleet {
		if courier.ZoneId != 0 {
			hardfleet[int(courier.ZoneId)] = append(hardfleet[int(courier.ZoneId)], courier)
		} else {
			couriers = append(couriers, courier)
		}
	}

	result := FleetSplitAlgo(couriers, hardfleet, ordersByZones)

	fmt.Println(result)
}
