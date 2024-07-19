package main

import (
	"math"
	"net/http"
	_ "net/http/pprof"
)

func hiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func main() {
	go Start()

	http.HandleFunc("/", hiHandler)
	http.ListenAndServe(":3001", nil)
}

func Start() {
	ReuseMemory2()
}

func ReuseMemory2() {

	sCurr1 := [][]int{[]int{1, 2, 3}, []int{4, 5, 6}}
	for i := 0; i < 10000; i++ {
		newSol := make([][]int, len(sCurr1))
		for j, r := range sCurr1 {
			newSol[j] = make([]int, len(r))
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
			sCurr1 = newSol
		}
	}

	sCurr2 := [][]int{{1, 2, 3}, {4, 5, 6}}
	var newSol [][]int
	for i := 0; i < 10000; i++ {
		if newSol == nil || len(newSol) != len(sCurr2) {
			newSol = make([][]int, len(sCurr2))
			for j := range sCurr2 {
				newSol[j] = make([]int, len(sCurr2[j]))
			}
		}
		for j, r := range sCurr2 {
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
			sCurr2 = newSol
		}
	}

	sl := [3][][]int{[][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 0}}, [][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 0}}, [][]int{[]int{0, 0, 0}, []int{0, 0, 0}, []int{0, 0, 0}}}
	sCurr := [][]int{[]int{1, 2, 3}, []int{4, 5, 6}}
	idx := 0
	for i := 0; i < 10000; i++ {
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

		if i != 2 && i != 5 {
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
		}

	}
}

func ReuseMemory() {
	sCurr := [][]int{[]int{1, 2, 3}, []int{4, 5, 6}}
	for i := 0; i < 10000; i++ {
		newSol := make([][]int, len(sCurr))
		for j, r := range sCurr {
			newSol[j] = make([]int, len(r))
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

func primeNumbers(max int) []int {
	var primes []int

	for i := 2; i < max; i++ {
		isPrime := true

		for j := 2; j <= int(math.Sqrt(float64(i))); j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primes = append(primes, i)
		}
	}

	return primes
}
