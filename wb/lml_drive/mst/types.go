package mst

type Param struct {
	Idx    int
	Min    float64
	Max    float64
	ShkMin int
	ShkMax int
}

type Params struct {
	p []Param
}

func (p Params) GetLowerShk(dist float64) int {
	for _, pp := range p.p {
		if dist >= pp.Min && dist < pp.Max {
			return pp.ShkMin
		}
	}

	return -1
}

func (p Params) GetUpperDist(dist float64) float64 {
	for _, pp := range p.p {
		if dist >= pp.Min && dist < pp.Max {
			return pp.Max
		}
	}

	return -1
}

func (p Params) GetUpperShk(dist float64) int {
	for _, pp := range p.p {
		if dist >= pp.Min && dist < pp.Max {
			return pp.ShkMax
		}
	}

	return -1
}

func (p Params) GetIdx(dist float64) int {
	for _, pp := range p.p {
		if dist >= pp.Min && dist < pp.Max {
			return pp.Idx
		}
	}

	return -1
}

type Edge struct {
	from int
	to   int
	w    float64
}
