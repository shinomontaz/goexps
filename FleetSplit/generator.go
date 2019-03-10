package main

import (
	"math/rand"
)

func generateFleet(TotalNum, HardNum int, Zones []Zone, tags []*Tag) []*Courier {
	fleet := make([]*Courier, 0, TotalNum)

	for i := 1; i <= HardNum; i++ {
		tagsNow := rand.Perm(len(tags))
		tagsNum := rand.Intn(len(tags))

		toTags := make([]*Tag, 0, len(tags))
		if rand.Float64() < 0.2 {
			for j := 0; j < tagsNum; j++ {
				toTags = append(toTags, tags[tagsNow[j]])
			}
		}

		fleet = append(fleet, &Courier{
			Id:     int64(i),
			ZoneId: int64(1 + rand.Intn(len(Zones)-1)),
			Limits: &Limits{
				Weight: float64(150 + rand.Intn(3)*50),
			},
			Tags: toTags,
		})
	}

	for i := 1; i <= TotalNum-HardNum; i++ {
		tagsNow := rand.Perm(len(tags))
		tagsNum := rand.Intn(len(tags))

		toTags := make([]*Tag, 0, len(tags))
		if rand.Float64() < 0.2 {
			for j := 0; j < tagsNum; j++ {
				toTags = append(toTags, tags[tagsNow[j]])
			}
		}

		fleet = append(fleet, &Courier{
			Id: int64(HardNum + i),
			Limits: &Limits{
				Weight: float64(150 + rand.Intn(3)*50),
			},
			Tags: toTags,
		})
	}

	//
	return fleet
}
