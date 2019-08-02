package types

import (
	"math"
	"time"
)

type LatLng struct {
	Lat float64
	Lng float64
}

type Courier struct {
		ID     int `json:"id"`
		Name   string `json:"name"`
		Limits struct {
			Weight int    `json:"weight"`
			Volume int    `json:"volume"`
			Orders string `json:"orders"`
		} `json:"limits"`
		WorkingTime struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"workingTime"`
		IsCycle bool        `json:"isCycle"`
		Geom    string      `json:"geom"`
		Center  interface{} `json:"center"`
		ZoneID  string      `json:"zoneId"`
		Tags    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"tags"`
		Wave           string      `json:"wave"`
		PhoneNumber    string      `json:"phoneNumber"`
		Type           string      `json:"type"`
		WarehouseID    string      `json:"warehouseId"`
		ExtID          string      `json:"extId"`
		CourierType    interface{} `json:"courierType"`
		ContractNumber string      `json:"contractNumber"`
		OrganizationID string      `json:"organizationId"`
}

type Order struct {
	ID             string    `json:"id"`
	GUID           string    `json:"guid"`
	Day            time.Time `json:"day"`
	CreatedAt      time.Time `json:"createdAt"`
	Address        string    `json:"address"`
	Lat            float64   `json:"lat"`
	Long           float64   `json:"long"`
	State          string    `json:"state"`
	TimeStart      time.Time `json:"timeStart"`
	Weight         float64   `json:"weight"`
	Volume         float64   `json:"volume"`
	Tags           []int     `json:"tags"`
	Warehouse      string    `json:"warehouse"`
	Wave           string    `json:"wave"`
	Route          string    `json:"route"`
	Courier        string    `json:"courier"`
	TimeEnd        time.Time `json:"timeEnd"`
	AddressFixed   string    `json:"addressFixed"`
	ClearingWaveID string    `json:"clearingWaveId"`
	ZoneID         string    `json:"zoneId"`
	ErrorReason    string    `json:"error_reason"`
}

func (from *LatLng) Distance(to *LatLng) float64 {
	dLat := (from.Lat - to.Lat) * math.Pi / 180.0
	dLng := (from.Lng - to.Lng) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(to.Lat*math.Pi/180.0)*math.Cos(from.Lat*math.Pi/180.0)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := 3958.75 * c
	return dist * 1609
}

func (p *LatLng) Angle(center *LatLng) float64 {
	rads := math.Atan2(p.Lng-center.Lng, p.Lat-center.Lat)
	if rads < 0 {
		rads += 2 * math.Pi
	}
	return rads
}
