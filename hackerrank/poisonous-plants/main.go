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
 * Complete the 'poisonousPlants' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts INTEGER_ARRAY p as parameter.
 */

type stack struct {
	items []int32
}

func NewStack() stack {
	return stack{
		items: make([]int32, 0),
	}
}

func (s *stack) isempty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}

func (s *stack) push(el int32) {
	s.items = append(s.items, el)
}

func (s *stack) pop() int32 {
	el := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return el
}

func (s *stack) peek() int32 {
	return s.items[len(s.items)-1]
}

func poisonousPlants(p []int32) int32 {
	// Write your code here
	s := make([]int32, 0)
	var max int32
	last_left := p[0]
	var steps int32
	for _, e := range p {
		if e <= last_left {
			s = append(s, 0)
		} else {
			s = append(s, 1)
			max = e
		}
		last_left = e
	}

	for max > 0 && len(s) > 1 {
		max = 0
		p2 := make([]int32, 0)
		for i, e := range s {
			if e == 0 {
				p2 = append(p2, p[i])
			}
		}

		p = p2

		s = make([]int32, 0)

		last_left = p[0]
		for _, e := range p {
			if e <= last_left {
				s = append(s, 0)
			} else {
				s = append(s, 1)
				max = 1
			}
			last_left = e
		}
		steps++
	}

	return steps
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	//	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	stdout, err := os.Create("out.txt")

	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	n := int32(nTemp)

	pTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	var p []int32

	for i := 0; i < int(n); i++ {
		pItemTemp, err := strconv.ParseInt(pTemp[i], 10, 64)
		checkError(err)
		pItem := int32(pItemTemp)
		p = append(p, pItem)
	}

	result := poisonousPlants(p)

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
