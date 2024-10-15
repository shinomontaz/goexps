package main

import (
	"lml-drive/types"
	"math"
)

const (
	earthEquator = 6378137.0
	earthExtr    = 1 / 298.257839
	earthPolus   = earthEquator * (1 - earthExtr)

	// osrmMultiplier -- поправочный к-т среднего роста расстояний при расчете по дорогам @see LGSW-1560.
	osrmMultiplier = 1.2809851063329336
)

// CalcSphereDistance -- считает сферическое расстояние между двумя точками.
// TODO требуется тестирование на точность вычислений. Расчет по среднему радиусу 6371км дает систематическую ошибку 0.5%.
func CalcSphereDistance(from, to types.Point) float64 {
	if from.Lat == to.Lat && from.Lng == to.Lng {
		return 0.0
	}

	fromLat := toRadian(from.Lat)
	toLat := toRadian(to.Lat)

	dlat := (fromLat - toLat) / 2
	dlng := toRadian(from.Lng-to.Lng) / 2
	sinLat := math.Sin(dlat)
	sinLng := math.Sin(dlng)

	A := sinLat*sinLat + math.Cos(fromLat)*math.Cos(toLat)*sinLng*sinLng
	distRad := 2.0 * math.Atan2(math.Sqrt(A), math.Sqrt(1.0-A))

	// !!! Учитываем эллипсоидность Земли на средней широте точек расчета:
	return earthAverage(fromLat, toLat) * distRad * osrmMultiplier
}

func toRadian(angle float64) float64 {
	return angle * math.Pi / 180.0
}

// earthAverage -- радиус на средней широте точек (задано в радианах).
func earthAverage(lat1, lat2 float64) float64 {
	return earthExtr*math.Cos((lat1+lat2)/2.0) + earthPolus
}
