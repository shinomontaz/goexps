package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

/*
 * Complete the 'componentsInGraph' function below.
 *
 * The function is expected to return an INTEGER_ARRAY.
 * The function accepts 2D_INTEGER_ARRAY gb as parameter.
 */

func componentsInGraph(gb [][]int32) []int32 {
	// Write your code here
	parent := make(map[int32]int32)
	groups := make(map[int32]map[int32]struct{})
	res := make([]int32, 2)

	for _, e := range gb {
		parent[e[0]] = e[0]
		if _, ok := parent[e[1]]; ok && parent[e[1]] != e[0] {
			// объединить группы
			parent[e[0]] = parent[e[1]]
		} else {
			parent[e[1]] = e[0]
		}
	}

	var find func(int32) int32
	find = func(x int32) int32 {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	for _, e := range gb {
		rootA := find(e[0])
		if _, ok := groups[rootA]; !ok {
			groups[rootA] = make(map[int32]struct{})
		}
		groups[rootA][e[0]] = struct{}{}
		rootB := find(e[1])
		if _, ok := groups[rootB]; !ok {
			groups[rootB] = make(map[int32]struct{})
		}
		groups[rootB][e[1]] = struct{}{}

	}

	for _, v := range groups {
		if res[0] == 0 || res[0] > int32(len(v)) {
			res[0] = int32(len(v))
		}
		if res[1] < int32(len(v)) {
			res[1] = int32(len(v))
		}
	}

	return res
}

func main() {
	// reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	// defer stdout.Close()

	// nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	// checkError(err)
	// n := int32(nTemp)

	// var gb [][]int32
	// for i := 0; i < int(n); i++ {
	// 	gbRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

	// 	var gbRow []int32
	// 	for _, gbRowItem := range gbRowTemp {
	// 		gbItemTemp, err := strconv.ParseInt(gbRowItem, 10, 64)
	// 		checkError(err)
	// 		gbItem := int32(gbItemTemp)
	// 		gbRow = append(gbRow, gbItem)
	// 	}

	// 	if len(gbRow) != 2 {
	// 		panic("Bad input")
	// 	}

	// 	gb = append(gb, gbRow)
	// }

	//	gb := [][]int32{[]int32{1, 6}, []int32{2, 7}, []int32{3, 8}, []int32{4, 9}, []int32{2, 6}, []int32{3, 7}}
	gb := [][]int32{[]int32{1, 6}, []int32{3, 8}, []int32{2, 6}, []int32{3, 7}}

	result := componentsInGraph(gb)

	fmt.Println(result)
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
