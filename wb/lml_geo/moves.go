package main

import (
	"lml/rand"
)

func insertmove(routes [][]int) [][]int {
	rndroute1 := rand.Intn(len(routes))
	rndroute2 := rand.Intn(len(routes))

	if rndroute1 == rndroute2 {
		return routes
	}
	if len(routes[rndroute1]) <= 1 {
		routes = append(routes[:rndroute1], routes[rndroute1+1:]...)
		return routes
	}
	if len(routes[rndroute2]) <= 1 {
		routes = append(routes[:rndroute2], routes[rndroute2+1:]...)
		return routes
	}

	var rndit1, rndit2 int
	if routes[rndroute1][0] == 0 {
		rndit1 = 1 + rand.Intn(len(routes[rndroute1])-1)
	} else {
		rndit1 = rand.Intn(len(routes[rndroute1]))
	}
	if routes[rndroute2][0] == 0 {
		rndit2 = 1 + rand.Intn(len(routes[rndroute2])-1)
	} else {
		rndit2 = rand.Intn(len(routes[rndroute2]))
	}

	it1 := routes[rndroute1][rndit1]

	if rndit1 != len(routes[rndroute1])-1 {
		routes[rndroute1] = append(routes[rndroute1][:rndit1], routes[rndroute1][rndit1+1:]...)
	} else {
		routes[rndroute1] = routes[rndroute1][:rndit1]
	}
	if rndit2 != len(routes[rndroute2])-1 {
		routes[rndroute2] = append(routes[rndroute2][:rndit2], append([]int{it1}, routes[rndroute2][rndit2:]...)...)
	} else {
		routes[rndroute2] = append(routes[rndroute2], it1)
	}

	return routes
}

func swapmove(routes [][]int) [][]int {
	rndroute1 := rand.Intn(len(routes))
	rndroute2 := rand.Intn(len(routes))

	if rndroute1 == rndroute2 {
		return routes
	}

	if len(routes[rndroute1]) <= 1 {
		routes = append(routes[:rndroute1], routes[rndroute1+1:]...)
		return routes
	}
	if len(routes[rndroute2]) <= 1 {
		routes = append(routes[:rndroute2], routes[rndroute2+1:]...)
		return routes
	}

	rndit1 := 1 + rand.Intn(len(routes[rndroute1])-1)
	rndit2 := 1 + rand.Intn(len(routes[rndroute2])-1)

	if rndroute1 == rndroute2 {
		routes[rndroute1][rndit1], routes[rndroute1][rndit2] = routes[rndroute1][rndit2], routes[rndroute1][rndit1]
		return routes
	}

	it1 := routes[rndroute1][rndit1]
	it2 := routes[rndroute2][rndit2]

	routes[rndroute1][rndit1] = it2
	routes[rndroute2][rndit2] = it1

	return routes
}

func twooptmove(routes [][]int) [][]int {
	rndroute1 := rand.Intn(len(routes))

	if len(routes[rndroute1]) <= 1 {
		routes = append(routes[:rndroute1], routes[rndroute1+1:]...)
		return routes
	}

	var rndit1, rndit2 int
	if routes[rndroute1][0] == 0 {
		rndit1 = 1 + rand.Intn(len(routes[rndroute1])-1)
		rndit2 = 1 + rand.Intn(len(routes[rndroute1])-1)
	} else {
		rndit1 = rand.Intn(len(routes[rndroute1]))
		rndit2 = rand.Intn(len(routes[rndroute1]))
	}

	if rndit1 == rndit2 {
		return routes
	}

	if rndit2 < rndit1 {
		rndit1, rndit2 = rndit2, rndit1
	}

	for k := 0; k <= (rndit2-rndit1)/2; k++ {
		routes[rndroute1][rndit1+k], routes[rndroute1][rndit2-k] = routes[rndroute1][rndit2-k], routes[rndroute1][rndit1+k]
	}

	return routes
}
