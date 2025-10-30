package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'solve' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts INTEGER_ARRAY arr as parameter.
 */

type stack struct {
	items []int32
}

func (s *stack) insert(e int32) int {
	for idx, el := range s.items {
		if el == e {
			return idx
		}
	}

	s.items = append(s.items, e)
	return len(s.items) - 1
}

func solve(arr []int32) int32 {
	// Write your code here
	var aux stack

	var pairs int32

	for _, e := range arr {
		idx := aux.insert(e)
		//        fmt.Printf("aux: %v, idx: %d, pairs was: %d\n", aux, idx, pairs)
		if idx < len(aux.items)-1 {
			all_less := true
			for i := idx; i < len(aux.items) && all_less; i++ {
				if e < aux.items[i] {
					all_less = false
				}
			}

			if all_less {
				pairs++
			}
		}
	}

	return pairs * 2
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	arrCount, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)

	arrTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	var arr []int32

	for i := 0; i < int(arrCount); i++ {
		arrItemTemp, err := strconv.ParseInt(arrTemp[i], 10, 64)
		checkError(err)
		arrItem := int32(arrItemTemp)
		arr = append(arr, arrItem)
	}

	result := solve(arr)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
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
