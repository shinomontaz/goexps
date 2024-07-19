package main

type ByPolar struct {
	list []int
	pts  []LatLong
}

func (a ByPolar) Len() int      { return len(a.list) }
func (a ByPolar) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a ByPolar) Less(i, j int) bool {
	return isPolarLess(a.list[i], a.list[j], a.pts)
}

// a less b means: a lays to left of b or a closer to center than b
func isPolarLess(p, q int, pts []LatLong) bool {
	c := pts[0]
	a := pts[p]
	b := pts[q]

	aLat := (a.Lat - c.Lat)
	bLat := (b.Lat - c.Lat)
	aLong := (a.Lng - c.Lng)
	bLong := (b.Lng - c.Lng)
	if (aLat >= 0) && (bLat < 0) {
		return true
	}
	if (aLat < 0) && (bLat >= 0) {
		return false
	}
	if aLat == 0 && bLat == 0 {
		if (aLong >= 0) || (bLong >= 0) {
			return a.Lng > b.Lng
		}
		return a.Lng < b.Lng
	}

	// dot product (center -> a) x (center -> b)
	dot := aLat*bLong - bLat*aLong
	if dot < 0 {
		return true
	}
	if dot > 0 {
		return false
	}

	return aLat*aLat+aLong*aLong < bLat*bLat+bLong*bLong // 2 points on 1 ray, return closest
}
