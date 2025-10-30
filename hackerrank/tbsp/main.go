package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'tbsp' function below.
 *
 * The function accepts following parameters:
 *  1. DOUBLE c
 *  2. DOUBLE d
 *  3. 2D_INTEGER_ARRAY cities
 */

func tbsp(c float64, d float64, cities [][]int32) {
	// Write your code here

}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	firstMultipleInput := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	nTemp, err := strconv.ParseInt(firstMultipleInput[0], 10, 64)
	checkError(err)
	n := int32(nTemp)

	c, err := strconv.ParseFloat(firstMultipleInput[1], 64)
	checkError(err)

	d, err := strconv.ParseFloat(firstMultipleInput[2], 64)
	checkError(err)

	var cities [][]int32
	for i := 0; i < int(n); i++ {
		citiesRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

		var citiesRow []int32
		for _, citiesRowItem := range citiesRowTemp {
			citiesItemTemp, err := strconv.ParseInt(citiesRowItem, 10, 64)
			checkError(err)
			citiesItem := int32(citiesItemTemp)
			citiesRow = append(citiesRow, citiesItem)
		}

		if len(citiesRow) != 3 {
			panic("Bad input")
		}

		cities = append(cities, citiesRow)
	}

	tbsp(c, d, cities)
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
