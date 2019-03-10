package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/shinomontaz/ga"
)

func FleetSplitAlgo(freeFleet []*Courier, hardFleet map[int][]*Courier, ordersByZones map[int][]*Order) map[int]int {
	// вычислить потребности по зонам
	// вычислить емкости жесткого флота по зонам
	// уменьшить потребности на емкости жесткого флота (?) - похоже на скелет расписания, а значит не надо ничего уменьшать

	zones := make([]Zone, 0, len(ordersByZones))

	for zoneId, orders := range ordersByZones {
		newZone := Zone{
			ZoneId: zoneId,
		}
		tagDemands := make(map[string]*TagDemand)
		tagDemands[""] = &TagDemand{}

		for _, ord := range orders {
			if len(ord.Tags) > 0 {
				tagsCode := ""
				for _, tagId := range ord.Tags {
					tagsCode = fmt.Sprintf("%s,%d", tagsCode, tagId)
				}
				if _, ok := tagDemands[tagsCode]; !ok {
					tagDemands[tagsCode] = &TagDemand{
						Tags: ord.Tags,
					}
				}
				tagDemands[tagsCode].Weight += ord.Parameters.Weight
			} else {
				tagDemands[""].Weight += ord.Parameters.Weight
			}
			newZone.Volume += ord.Parameters.Volume
			newZone.Weight += ord.Parameters.Weight
			newZone.Demand = tagDemands
		}
		zones = append(zones, newZone)
	}

	PFactory := &SolutionFactory{
		Zones:     zones,
		core:      hardFleet,
		freeFleet: freeFleet,
	}
	var gaSolver = ga.Ga{
		NewIndividual: PFactory.NewIndividual,
		PopSize:       30,
	}

	gaSolver.Initialize()
	record := gaSolver.Population[0].Clone()

	bestFitness := record.Fitness()
	epochMax := 1000
	epoch := 1
	for epoch < epochMax {
		gaSolver.Evolve()
		epoch++
		curr := gaSolver.Record()
		if curr.Fitness() < bestFitness { // чем меньше, тем лучше
			bestFitness = curr.Fitness()
			record = curr.Clone()
			fmt.Println(bestFitness)
		}
	}

	result := make(map[int]int, len(record.(*Split).List))
	for zoneId, couriers := range record.(*Split).List {
		result[zoneId] = len(couriers)
	}

	fmt.Printf("Result: %v, Score: %f\n", record, record.Fitness())

	return result
}

type SolutionFactory struct {
	Zones     []Zone
	core      map[int][]*Courier // распределение жесткого флота
	freeFleet []*Courier
}

func (sf *SolutionFactory) NewIndividual() ga.Individual {
	nSplit := &Split{
		List:       sf.core,
		Couriers:   sf.freeFleet,
		CourierMap: make(map[int64]int),
	}

	for _, zone := range sf.Zones {
		nSplit.Zones = append(nSplit.Zones, zone)
	}

	for zoneId, couriers := range sf.core {
		for _, courier := range couriers {
			nSplit.CourierMap[courier.Id] = zoneId
		}
	}

	for _, courier := range sf.freeFleet {
		randZone := rand.Intn(len(sf.Zones))
		nSplit.AddCourier(courier, sf.Zones[randZone].ZoneId)
		nSplit.CourierMap[courier.Id] = sf.Zones[randZone].ZoneId
	}
	return nSplit
}

func (s *Split) Mutate() ga.Individual {
	cl := s.Clone()

	// взять случайную зону и случайного курьера из свободных (!)
	// перекинуть    return cl

	return cl
}

func (s *Split) Fitness() float64 {
	unfitness := 0.0
	for _, zone := range s.Zones {
		// vSupply := 0.0
		// wSupply := 0.0

		tagSupply := make(map[string]*TagDemand)
		tagSupply[""] = &TagDemand{}

		for _, cou := range s.List[zone.ZoneId] {

			if len(cou.Tags) > 0 {
				tagsCode := ""
				couTags := make([]int, 0, len(cou.Tags))
				for _, tag := range cou.Tags {
					tagsCode = fmt.Sprintf("%s,%d", tagsCode, tag.Id)
					couTags = append(couTags, tag.Id)
				}
				if _, ok := tagSupply[tagsCode]; !ok {
					tagSupply[tagsCode] = &TagDemand{Tags: couTags}
				}
				tagSupply[tagsCode].Weight += cou.Limits.Weight
			} else {
				tagSupply[""].Weight += cou.Limits.Weight
			}

			// vSupply += cou.Limits.Volume
			// wSupply += cou.Limits.Weight
		}

		unfitness += math.Max(math.Max(zone.Weight-wSupply, zone.Volume-vSupply), 0)
	}
	return 1 / (unfitness + 1.0)
}

func (s *Split) Clone() ga.Individual {
	cl := &Split{
		List:       make(map[int][]*Courier),
		Couriers:   s.Couriers,
		CourierMap: make(map[int64]int),
		Zones:      s.Zones,
	}

	for couId, zoneId := range s.CourierMap {
		cl.CourierMap[couId] = zoneId
	}

	for zoneId, couList := range s.List {
		cl.List[zoneId] = make([]*Courier, 0, len(couList))
		for _, cou := range couList {
			cl.List[zoneId] = append(cl.List[zoneId], cou)
		}
	}
	return cl
}

func (s *Split) Educate() {}

func (s *Split) Crossover(p ga.Individual) ga.Individual {
	child := s.Clone()
	for _, courier := range s.Couriers {
		courierZone := 0
		if _, ok1 := s.CourierMap[courier.Id]; ok1 {
			if _, ok2 := p.(*Split).CourierMap[courier.Id]; ok2 {
				if s.CourierMap[courier.Id] == p.(*Split).CourierMap[courier.Id] {
					courierZone = s.CourierMap[courier.Id]
				}
			}
		}
		if courierZone == 0 {
			courierZone = s.Zones[rand.Intn(len(s.Zones))].ZoneId
		}
		child.(*Split).RemoveCourier(courier)
		child.(*Split).AddCourier(courier, courierZone)
	}
	return child
}

func (s *Split) AddCourier(cou *Courier, zoneId int) {
	s.List[zoneId] = append(s.List[zoneId], cou)
}

func (s *Split) RemoveCourier(cou *Courier) {
	for zoneId, couList := range s.List {
		for i, cour := range couList {
			if cour.Id == cou.Id {
				s.List[zoneId] = append(s.List[zoneId][:i], s.List[zoneId][i+1:]...)
				return
			}
		}
	}
}
