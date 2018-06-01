package model

import (
	"fmt"
	"log"
)

type Zone struct {
	ID      int32
	Name    string
	GeoJson string
	Fleet   []*Courier
}

var listZones map[int32]*Zone

func GetFleet() (zones map[int32]*Zone) {
	zones = make(map[int32]*Zone)
	sql := `SELECT z.id, z.name, f.id, f.name, f.volume, f.weight FROM zone z RIGHT JOIN fleet f ON f.fk_zone = z.id`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var zID, fID int32
	var zName, fName, zGeojson string
	var fWeight, fVolume float64

	for rows.Next() {
		err = rows.Scan(&zID, &zName, &fID, &fName, &fVolume, &fWeight)
		if err != nil {
			log.Fatal(err)
		}

		if _, exists := zones[zID]; !exists {
			zones[zID] = &Zone{ID: zID, Name: zName, GeoJson: zGeojson, Fleet: make([]*Courier, 0, 1)}
		}
		zones[zID].Fleet = append(zones[zID].Fleet, &Courier{Id: fID, Name: fName, Volume: fVolume, Weight: fWeight})
	}

	return zones
}

func GetZones() (zones map[int32]*Zone) {
	if len(listZones) > 0 {
		return listZones
	}

	listZones = make(map[int32]*Zone)
	sql := `SELECT z.id, z.name,
		jsonb_build_object(
        'type',       'Feature',
        'id',         z.id,
        'geometry',   ST_AsGeoJSON(z.geom)::jsonb,
        'properties', to_jsonb(z) - 'id' - 'geom'
		) AS geojson FROM zone z`

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
			listZones[zID] = &Zone{ID: zID, Name: zName, GeoJson: zGeojson}
		}
	}

	return listZones
}

func GetPointZone(lat, long float64) *Zone {
	zones := GetZones()
	for _, zone := range zones {
		if zone.IsPointInside(lat, long) {
			return zone
		}
	}

	return nil
}

func (z *Zone) IsPointInside(lat, long float64) (contains bool) {
	sql := `SELECT (
		ST_Contains( (SELECT geom FROM zone WHERE id = %d ),
		  ST_GeomFromText('POINT(%f %f)', 4326) )
		OR
		ST_Contains( (SELECT ST_Makeline( ST_Boundary(geom) ) FROM zone WHERE id = %d ),
		  ST_GeomFromText('POINT(%f %f)', 4326) )
	  )  AS contains;`

	sql = fmt.Sprintf(sql, z.ID, lat, long, z.ID, lat, long)

	db.QueryRow(sql).Scan(&contains)

	return contains
}
