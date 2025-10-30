package main

type LatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Options struct {
	CourierCapacity int `json:"courier_capacity"`
	UnloadDuration  int `json:"unload_duration"`
	LoadDuration    int `json:"load_duration"`
}

type Point struct {
	Id    int64  `json:"id"`
	Coord LatLon `json:"coord"` // Точка назначения
	Shk   int    `json:"shk"`   // Количество штук в заказе
}

type TargetPoint struct {
	PointId    int64   `json:"PointId"`
	SupplierID int64   `json:"SupplierId"`
	OrderID    int64   `json:"OrderId"`
	Route      int     `json:"Route"`
	Num        int     `json:"Num"`
	Tw         string  `json:"Tw"`
	Lat        float64 `json:"Lat"`
	Lon        float64 `json:"Lon"`
}

type TargetRoute struct {
	Points   []int64 `json:"Points"`
	Volume   int     `json:"Volume"`
	Time     int     `json:"Time"`
	Distance float64 `json:"Distance"`
}

type TargetData struct {
	Pts           []TargetPoint `json:"Pts"`
	Routes        []TargetRoute `json:"Routes"`
	TotalTime     int           `json:"TotalTime"`
	TotalDistance int           `json:"TotalDistance"`
}
