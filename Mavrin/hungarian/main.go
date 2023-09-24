package main

import (
	"fmt"
	"math"
)

type testCase struct {
	in   [][]int
	want []int
}

func main() {

	cases := []testCase{
		{
			in: [][]int{
				{11, 6, 12},
				{12, 4, 6},
				{8, 12, 11},
			},
			want: []int{1, 2, 0},
		},
		{
			in: [][]int{
				{13, 13, 19, 50, 33, 38},
				{73, 33, 71, 77, 97, 95},
				{20, 8, 56, 55, 64, 35},
				{26, 25, 72, 32, 55, 77},
				{83, 40, 69, 3, 53, 49},
				{67, 20, 44, 29, 86, 61},
			},
			want: []int{4, 1, 5, 0, 3, 2},
		},
	}

	for _, c := range cases {
		hungarian(c.in)
		fmt.Println(c.want)
	}
}

func hungarian(c [][]int) []int {
	// 1. check matrix to proove solution exists
	// 2. start algo:
	// 2.0. find perfect matching on 0-cell basis. If not possible:
	// 	2.1. -min for all rows ( use potentials )
	// 	2.2. find submatrix and make a -min for all rows (update L- ) and +min for not used cols ( update R+ ) - use a DFS to find
	// 	2.3. find perfect matching

	// OR

	// for v in L
	// while true
	// 	if dfs(v)
	// 		M++
	// 		break
	// 	else
	// 		DELTA = min Cuv, u in L+, v in R-
	// 		Cuv -= DELTA | u in L+
	// 		Cuv += DELTA | v in R+

	// for v in L
	// 	// run DFS from v:
	// 	v in L+, others in L-
	// 	R- = R
	// 	m[j] = Cvj // minimum of a j col
	// 	// trying to find new edge from L+ to R-
	// 	while true:
	// 		D = min(m[j]) | j in R- // for ... , also find an edge (u, v) where min occurs
	// 		// recalc potentials: a is a row potential, b - is a col potential
	// 		a[i] -= D | i in L+
	// 		b[j] += D | j in R+
	// 		m[j] -= D
	// 		R+ = R+U{u}
	// 		if p(u) == nil // p(u) - pair vertex of u
	// 			// inverse path
	// 			break
	// 		else
	// 			L+ = L+U{p(u)}
	// 			// update mins:
	// 			m[j] = min(m[j], C[p(u), j] + a[p(u)] + b[j])

	inf := 0 // max element
	for _, row := range c {
		for _, el := range row {
			if inf < el {
				inf = el
			}
		}
	}

	n := len(c)
	pr := make([]int, n+1)
	pl := make([]int, n+1)
	mt := make([]int, n+1)
	for i := 0; i <= n; i++ {
		mt[i] = -1
	}
	dist := make([]int, n+1)
	par := make([]int, n+1)
	used := make([]bool, n+1)

	for row := 1; row <= n; row++ {
		mt[0] = row
		used = make([]bool, n+1)
		for i := 1; i <= n; i++ {
			dist[i] = inf // ?
		}

		col := 0
		for mt[col] != -1 {
			u := mt[col]
			used[col] = true
			for i := 1; i <= n; i++ {
				if !used[i] {
					x := c[u-1][i-1] - pl[u-1] - pr[i-1]
					if x < dist[i] {
						dist[i] = x
						par[i] = col
					}
				}
			}

			delta := inf // ?

			for i := 1; i <= n; i++ {
				if !used[i] {
					delta = int(math.Min(float64(delta), float64(dist[i])))
				}
			}

			for i := 0; i <= n; i++ { // here from 0
				if used[i] {
					pr[i] -= delta
					pl[mt[i]] += delta
				}
			}

			for i := 1; i <= n; i++ {
				if !used[i] {
					dist[i] -= delta
				}
				if dist[i] == 0 {
					col = i
				}
			}
		}

		for col != 0 {
			mt[col] = mt[par[col]]
			col = par[col]
		}
	}

	fmt.Println(par)

	return par
}
