package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

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
					if tagsCode != "" {
						tagsCode = fmt.Sprintf("%s,", tagsCode)
					}
					tagsCode = fmt.Sprintf("%s%d", tagsCode, tagId)
				}
				if _, ok := tagDemands[tagsCode]; !ok {
					tagDemands[tagsCode] = &TagDemand{
						Tags: ord.Tags,
						Key:  tagsCode,
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

	fmt.Println(bestFitness)

	epochMax := 100
	epoch := 1
	for epoch < epochMax {
		gaSolver.Evolve()
		epoch++
		curr := gaSolver.Record()
		if curr.Fitness() > bestFitness { // чем меньше, тем лучше
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
		vSupply := 0.0
		wSupply := 0.0

		tagSupply := make(map[string]*TagDemand)
		tagSupply[""] = &TagDemand{}

		for _, cou := range s.List[zone.ZoneId] {

			if len(cou.Tags) > 0 {
				tagsCode := ""
				couTags := make([]int, 0, len(cou.Tags))
				for _, tag := range cou.Tags {
					if tagsCode != "" {
						tagsCode = fmt.Sprintf("%s,", tagsCode)
					}
					tagsCode = fmt.Sprintf("%s%d", tagsCode, tag.Id)

					//					tagsCode = fmt.Sprintf("%s,%d", tagsCode, tag.Id)

					couTags = append(couTags, tag.Id)
				}
				if _, ok := tagSupply[tagsCode]; !ok {
					tagSupply[tagsCode] = &TagDemand{
						Tags: couTags,
						Key:  tagsCode,
					}
				}
				tagSupply[tagsCode].Weight += cou.Limits.Weight
			} else {
				tagSupply[""].Weight += cou.Limits.Weight
			}

			// vSupply += cou.Limits.Volume
			// wSupply += cou.Limits.Weight
		}

		// 1. берем требования с максимальной мощностью по тегам
		// 2. ищем такой же или охватывающий набор тегов в поддерживаемых для рассматриваемого
		// 2.1 нашли - снижаем емкость соответствующего поддерживаемого набора
		// 2.2 не нашли - увеличиваем unfitness и пишем в требования

		for {
			biggestDemand := zone.TakeBiggestDemand()
			if biggestDemand == nil {
				break
			}
			equalOrBiggerSupply := findEqualOrBigger(biggestDemand, tagSupply)
			if equalOrBiggerSupply == nil {
			} else {
				wSupply = equalOrBiggerSupply.Weight
				vSupply = equalOrBiggerSupply.Volume
				equalOrBiggerSupply.Weight = math.Max(wSupply-biggestDemand.Weight, 0)
				equalOrBiggerSupply.Volume = math.Max(vSupply-biggestDemand.Volume, 0)

				if equalOrBiggerSupply.Weight == 0 && equalOrBiggerSupply.Volume == 0 {
					delete(tagSupply, equalOrBiggerSupply.Key)
				}
			}

			if biggestDemand.Key == "" && math.Max(math.Max(biggestDemand.Weight-wSupply, biggestDemand.Volume-vSupply), 0) > 0 {
				for equalOrBiggerSupply1 := findEqualOrBigger(biggestDemand, tagSupply); equalOrBiggerSupply1 != nil; {
					wSupply += equalOrBiggerSupply1.Weight
					vSupply += equalOrBiggerSupply1.Volume
					equalOrBiggerSupply1.Weight = math.Max(wSupply-biggestDemand.Weight, 0)
					equalOrBiggerSupply1.Volume = math.Max(vSupply-biggestDemand.Volume, 0)

					if equalOrBiggerSupply1.Weight == 0 && equalOrBiggerSupply1.Volume == 0 {
						delete(tagSupply, equalOrBiggerSupply1.Key)
					}
				}
			}

			unfitness += math.Max(math.Max(biggestDemand.Weight-wSupply, biggestDemand.Volume-vSupply), 0)

		}

	}
	return 1 / (unfitness + 1.0)
}

func findEqualOrBigger(tgd *TagDemand, tgSupply map[string]*TagDemand) *TagDemand {
	if tgd.Key == "" {
		if _, ok := tgSupply[""]; ok {
			return tgSupply[""]
		}
		return nil
	}

	tgSupplySl := make([]*TagDemand, 0, len(tgSupply))

	for _, tgs := range tgSupply {
		tgSupplySl = append(tgSupplySl, tgs)
	}

	sort.Slice(tgSupplySl, func(i, j int) bool { // здесь сперва будут те, у кого количество тегов меньше
		return len(tgSupplySl[i].Tags) < len(tgSupplySl[j].Tags)
	})

	for _, supply := range tgSupplySl { // от меньших к большему
		if len(tgd.Key) > len(supply.Key) {
			continue
		}
		if tgd.Key == supply.Key {
			return tgSupply[supply.Key]
		}
		if strings.Index(supply.Key, tgd.Key) != -1 {
			return tgSupply[supply.Key]
		}
	}

	return nil
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
