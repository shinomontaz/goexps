package main

import "sort"

type Tag struct {
	Id   int
	Name string
	Code string
}

type Limits struct {
	Weight float64
	Volume float64
	Orders int64
}

type Courier struct {
	Id     int64
	Name   string
	Limits *Limits
	ZoneId int64
	Tags   []*Tag
}

type Split struct {
	List       map[int][]*Courier
	Zones      []Zone
	CourierMap map[int64]int
	Couriers   []*Courier
}

type LatLong struct {
	Lat  float64
	Long float64
}

type Order struct {
	Coords      *LatLong
	Address     string
	Guid        string
	Parameters  *Limits
	ServiceTime int64
	Zone        int64
	Tags        []int
}

type Zone struct {
	Volume float64
	Weight float64
	ZoneId int
	Demand map[string]*TagDemand
}

type TagDemand struct {
	Tags   []int
	Weight float64
	Volume float64
	Key    string
}

func (z *Zone) TakeBiggestDemand() *TagDemand {
	if len(z.Demand) == 0 {
		return nil
		//		panic("zone cannot have empty demand! at least we should have a demand without tags?")
	}

	tgDemands := make([]*TagDemand, 0, len(z.Demand))

	for _, tgd := range z.Demand {
		tgDemands = append(tgDemands, tgd)
	}

	sort.Slice(tgDemands, func(i, j int) bool {
		return len(tgDemands[i].Tags) > len(tgDemands[j].Tags)
	})

	res := z.Demand[tgDemands[0].Key]
	delete(z.Demand, tgDemands[0].Key)
	return res
}
