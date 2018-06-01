package model

import (
	"log"
)

type Zone struct {
	ID      int32
	Name    string
	GeoJson string
	Fleet   []*Courier
}

func GetFleet() (zones map[int32]*Zone) {

	sql := `SELECT z.id, z.name, f.id, f.name, f.volume, f.weight FROM zone RIGHT JOIN fleet f ON f.fk_zone = z.id`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var zID, fID int32
	var zName, fName, zGeojson string
	var fWeight, fVolume float64

	for rows.Next() {
		err = rows.Scan(&zID, &zName, &zGeojson, &fID, &fName, &fVolume, &fWeight)
		if err != nil {
			log.Fatal(err)
		}

		if _, exists := zones[zID]; !exists {
			zones[zID] = &Zone{ID: zID, Name: zName, GeoJson: zGeojson, Fleet: make([]*Courier, 1)}
		}
		zones[zID].Fleet = append(zones[zID].Fleet, &Courier{Id: fID, Name: fName, Volume: fVolume, Weight: fWeight})
	}

	return zones
}

func GetZones() (zones map[int32]*Zone) {

	sql := `SELECT z.id, z.name,
		jsonb_build_object(
        'type',       'Feature',
        'id',         z.id,
        'geometry',   ST_AsGeoJSON(z.geom)::jsonb,
        'properties', to_jsonb(row) - 'id' - 'geom'
		) AS geojson`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var zID int32
	var zName, zGeojson string

	for rows.Next() {
		err = rows.Scan(&zID, &zName, &zGeojson)
		if err != nil {
			log.Fatal(err)
		}

		if _, exists := zones[zID]; !exists {
			zones[zID] = &Zone{ID: zID, Name: zName, GeoJson: zGeojson}
		}
	}

	return zones
}
