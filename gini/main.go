package main

import (
	"fmt"

	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
)

func main() {
	g := gini.New()

	n := 8
	board := make([][]z.Lit, n)

	for i := 0; i < n; i++ {
		board[i] = make([]z.Lit, n)
		for j := 0; j < n; j++ {
			board[i][j] = z.Var(n + i*n + j).Pos() // Магия gini - если подать 0 на Var, то будет неправильно ибо 0 тут трактуется как безусловная ложь
		}
	}

	// как минимум одна королева в столбце
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			g.Add(board[i][j])
		}
		g.Add(0)
	}

	// как минимум одна королева в строке
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			g.Add(board[j][i])
		}
		g.Add(0)
	}

	// максимум одна королева в строке
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := j + 1; k < n; k++ {
				g.Add(board[i][j].Not())
				g.Add(board[i][k].Not())
				g.Add(0)
			}
		}
	}

	// максимум одна королева в столбце
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := j + 1; k < n; k++ {
				g.Add(board[j][i].Not())
				g.Add(board[k][i].Not())
				g.Add(0)
			}
		}
	}

	// главная диагональ
	for i := n - 1; i > 0; i-- {
		for j := 0; j < n; j++ {
			for ii := i - 1; ii > 0; ii-- {
				for jj := j + 1; jj < n; jj++ {
					if i+j == ii+jj {
						g.Add(board[i][j].Not())
						g.Add(board[ii][jj].Not())
						g.Add(0)
					}
				}
			}
		}
	}

	// побочная диагональ
	for i := 0; i < n; i++ {
		for j := n - 1; j > 0; j-- {
			for ii := i + 1; ii < n; ii++ {
				for jj := j - 1; jj > 0; jj-- {
					if i-j == ii-jj {
						g.Add(board[i][j].Not())
						g.Add(board[ii][jj].Not())
						g.Add(0)
					}
				}
			}
		}
	}

	if g.Solve() != 1 {
		fmt.Printf("error.\n")
		for i := 0; i < n; i++ {
			fmt.Println(g.Why(board[i]))
			break
		}
		return
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if g.Value(board[i][j]) {
				fmt.Printf("*")
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
