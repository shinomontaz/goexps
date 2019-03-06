package main

type Tag struct {
	Id   int64
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
	List map[int64][]*Courier
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
}

type Zone struct {
	Volume float64
	Weight float64
	ZoneId int64
}
