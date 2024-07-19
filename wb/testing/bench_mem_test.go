package main

import (
	"testing"
)

func copy2DArray(src [][]int) [][]int {
	dst := make([][]int, len(src))
	for i, row := range src {
		dst[i] = make([]int, len(row))
		copy(dst[i], row)
	}
	return dst
}

func benchmarkCopyFunction(b *testing.B) {
	sCurr := [][]int{{1, 2, 3}, {4, 5, 6}}
	for i := 0; i < b.N; i++ {
		newSol := copy2DArray(sCurr)

		if i == 3 {
			newSol = append(newSol, []int{7, 8, 9})
		}
		if i == 6 {
			newSol = append(newSol[:1], newSol[2:]...)
		}

		for _, r := range newSol {
			r[0]++
		}

		if i != 2 && i != 5 {
			sCurr = newSol
		}
	}
}

func BenchmarkCopyFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkCopyFunction(b)
	}
}

func benchmarkReuseMemory(b *testing.B) {
	sCurr := [][]int{{1, 2, 3}, {4, 5, 6}}
	var newSol [][]int
	for i := 0; i < b.N; i++ {
		if newSol == nil || len(newSol) != len(sCurr) {
			newSol = make([][]int, len(sCurr))
			for j := range sCurr {
				newSol[j] = make([]int, len(sCurr[j]))
			}
		}
		for j, r := range sCurr {
			copy(newSol[j], r)
		}

		if i == 3 {
			newSol = append(newSol, []int{7, 8, 9})
		}
		if i == 6 {
			newSol = append(newSol[:1], newSol[2:]...)
		}

		for _, r := range newSol {
			r[0]++
		}

		if i != 2 && i != 5 {
			sCurr = newSol
		}
	}
}

func benchmarkReuseMemory2(b *testing.B) {
	sl := [3][][]int{[][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 0}}, [][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 0}}, [][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 0}}}
	sCurr := [][]int{[]int{1, 2, 3}, []int{4, 5, 6}}
	idx := 0
	for i := 0; i < b.N; i++ {
		idx = idx % len(sl)
		newSol := sl[idx]
		newSol = newSol[:len(sCurr)]
		for j, r := range sCurr {
			newSol[j] = newSol[j][:len(r)]
			copy(newSol[j], r)
		}
		if i == 3 {
			newSol = append(newSol, []int{7, 8, 9})
		}

		if i == 6 {
			newSol = append(newSol[:1], newSol[2:]...)
		}

		for _, r := range newSol {
			r[0]++
		}

		//		if i != 2 && i != 5 {
		if len(sCurr) < len(newSol) {
			sCurr = append(sCurr, make([][]int, len(newSol)-len(sCurr))...)
		}
		for j, r := range newSol {
			if len(sCurr[j]) < len(r) {
				sCurr[j] = append(sCurr[j], make([]int, len(r)-len(sCurr[j]))...)
			}
			sCurr[j] = sCurr[j][:len(r)]
			copy(sCurr[j], r)
		}

		sCurr = sCurr[:len(newSol)]
		idx++
		//		}

	}
}

func BenchmarkReuseMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkReuseMemory(b)
	}
}

func BenchmarkReuseMemory2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkReuseMemory2(b)
	}
}
