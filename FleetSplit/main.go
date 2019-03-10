package main

import (
	"fmt"
	"math/rand"
	"sort"
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
	ordersNum := 5000
	zones := make([]Zone, 0, zonesNum)
	orders := make([]*Order, 0, ordersNum)
	for i := 0; i < ordersNum; i++ {
		tagsNow := rand.Perm(len(tags))
		tagsNum := rand.Intn(len(tags))

		toTags := make([]int, 0, len(tags))
		if rand.Float64() < 0.1 {
			for j := 0; j < tagsNum; j++ {
				toTags = append(toTags, tags[tagsNow[j]].Id)
			}
		}

		sort.Ints(toTags)
		orders = append(orders, &Order{
			Zone: int64(rand.Intn(ordersNum) % zonesNum),
			Tags: toTags,
			Parameters: &Limits{
				Weight: float64(rand.Intn(5) + (rand.Intn(10)%5)*5),
			},
		})
	}

	ordersByZones := make(map[int][]*Order, 0)
	for _, ord := range orders {
		ordersByZones[int(ord.Zone)] = append(ordersByZones[int(ord.Zone)], ord)
	}

	// generate zones
	for i := 0; i < zonesNum; i++ {
		newZone := Zone{
			ZoneId: i,
		}
		zones = append(zones, newZone)
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
