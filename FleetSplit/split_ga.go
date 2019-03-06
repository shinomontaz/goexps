package main

import (
	"math/rand"

	"github.com/shinomontaz/ga"
)

func FleetSplitAlgo(freeFleet []*Courier, hardFleet map[int][]*Courier, ordersByZones map[int][]*Order, courierNums map[int]int) {
	// вычислить потребности по зонам
	// вычислить емкости жесткого флота по зонам
	// уменьшить потребности на емкости жесткого флота (?) - похоже на скелет расписания, а значит не надо ничего уменьшать
	//
}

type SolutionFactory struct {
	Zones     []Zone
	core      map[int64][]*Courier // распределение жесткого флота
	freeFleet []*Courier
}

func (sf *SolutionFactory) GenerateSolution() ga.Individual {
	nSplit := &Split{List: sf.core}
	for _, courier := range sf.freeFleet {
		randZone := rand.Intn(len(sf.Zones))
		nSplit.AddCourier(courier, sf.Zones[randZone].ZoneId)
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
	return 0.0
}

func (s *Split) Crossover(partner ga.Individual) ga.Individual {
	cl := s.Clone()
	return cl
}

func (s *Split) Clone() ga.Individual {
	cl := &Split{List: make(map[int64][]*Courier)}
	for zoneId, couList := range s.List {

		cl.List[zoneId] = make([]*Courier, 0, len(couList))

		for _, cou := range couList {

			cl.List[zoneId] = append(cl.List[zoneId], cou)

		}

	}
	return cl
}

func (s *Split) Educate() {}

func (s *Split) AddCourier(cou *Courier, zoneId int64) {
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
