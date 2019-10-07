package main

import (
	"fmt"

	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
)

func main() {
	g := gini.New()

	n := 8

	board := make([]z.Lit, n*n)
	// var board = func(i, j int, present bool) z.Lit {
	// 	return z.Var(i*n + j).Pos()
	// }
	m := []z.Lit{z.Var(1).Pos(), z.Var(2).Pos(), z.Var(3).Pos()}

	g.Add(m[0])
	g.Add(m[1])
	g.Add(m[2])
	g.Add(0)

	for j := 0; j < 2; j++ {
		for k := j + 1; k <= 2; k++ {
			g.Add(m[j].Not())
			g.Add(m[k].Not())
			g.Add(0)
		}
	}

	if g.Solve() == 1 {
		for j := 0; j <= 2; j++ {
			fmt.Printf("%d (%d)\n", 1, g.Value(m[j]))
		}
		//		panic("active!")
		// do something
	}
	fmt.Println(g.Why(m))

	panic("false!")

	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			board[(i-1)*n+(j-1)] = z.Lit((i-1)*n + (j - 1)) // z.Var((i-1)*n + (j - 1)).Pos() // z.Var() //
			fmt.Println(i, j, n, i*n+j, board[(i-1)*n+(j-1)])
			g.Add(board[(i-1)*n+(j-1)])
		}
		g.Add(0) // gini magic here
	}

	if g.Solve() != 1 {
		fmt.Printf("error.\n")
		return
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if g.Value(board[(i-1)*n+(j-1)]) {
				fmt.Printf("%d (%d)", 1, g.Value(board[(i-1)*n+(j-1)]))
				//				break
			} else {
				fmt.Printf("%d", 0)
			}
			if j != n {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}
