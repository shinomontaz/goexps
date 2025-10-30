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
