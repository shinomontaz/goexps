package main

import (
	"fmt"
	"math"
)

// Функция должна вернуть значение number, где бит на позиции bitIndex сброшен в 0
func resetBit(number, bitIndex int) int {
	s := make([]int, 0)
	for number > 0 {
		s = append(s, number%2)
		number /= 2
	}
	if len(s) == 0 {
		s = append(s, 0)
	}

	ss := make([]int, 0)
	for i := len(s) - 1; i >= 0; i-- {
		ss = append(ss, s[i])
	}

	fmt.Println(ss)
	ss[bitIndex] = 0
	fmt.Println(ss)
	res := 0

	for i := 0; i < len(ss); i++ {
		if ss[i] > 0 {
			res += int(math.Pow(2, float64(len(ss)-i-1)))
		}
	}

	return res

	// КОДИТЬ СЮДА
}

func main() {
	fmt.Printf("0, 0 => %d (expected 0)\n", resetBit(0, 0))
	fmt.Printf("3, 0 => %d (expected 2)\n", resetBit(3, 0))
	fmt.Printf("3, 1 => %d (expected 1)\n", resetBit(3, 1))
	fmt.Printf("3, 2 => %d (expected 3)\n", resetBit(16, 0))
}
