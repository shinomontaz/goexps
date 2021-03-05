package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const MOD = 1000000009
const ALPHABET_LEN = 26
const MAX_K = 26
const MAX_N = 100000

var precalc [][]int32

/*
 * Complete the dortmundDilemma function below.
 */
func dortmundDilemma(n int32, k int32) int32 {
	/*
	 * Write your code here.
	 */

	if k > MAX_K || k >= n || k > (n/2+n%2) {
		return 0
	}

	return int32(solve(n, k)*Combinations(ALPHABET_LEN, k)) % MOD

	// dynamic part

}

// presolve for n=2k
// for 2k-2 < n < 2k reduce to (n=2k-2, k=k-1)

func solve(n, k int32) int32 {
	// here n > 1

	if k == 1 {
		return 1
	}

	if k >= n || k > (n/2+n%2) {
		return 0
	}

//	if n == 2*k {

//		return precalc[k][n]
//	}

	if n%2 > 0 {
		return solve(n-1, k) * k
	}
	if n > 2*k {
		return solve(n-2, k) * k * k
	}

	res := 0
	for i := 0; i<=k-1; i++ {
		res += Permutations(k-i, k-i) * (math.Pow(k, 2 * i) - solve( k, 2 * i )  )
	}

	return res
}

func Permutations(n, k int32) int64 {
	if k > n {
		panic("k > n!")
	}

	var i int32
	var res int64
	res = 1
	for i = 0; i <= k-1; i++ {
		res *= int64(n - i)
	}
	return res
}

func Combinations(n, k int32) int64 {
	if k > n {
		panic("k > n!")
	}

	var i int32
	var res int64
	res = 1
	for i = 0; i <= k-1; i++ {
		res *= int64(n - i)
	}
	kfact := 1
	for i = 1; i <= k; i++ {
		kfact *= i
	}

	res /= kfact

	return res
}


func precalcuations() {
	precalc = make([][]int32, MAX_K)
	for k = 1; k <= MAX_K; k++ {
		n = 2*k
		precalc[k][n] = 
	}
}

func main() {

	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout := os.Stdout

	// stdout, err := os.Stdout
	// Create(os.Getenv("OUTPUT_PATH"))
	// checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	tTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	t := int32(tTemp)

	// recalculations for n=2k
	precalcuations()

	for tItr := 0; tItr < int(t); tItr++ {
		nk := strings.Split(readLine(reader), " ")

		nTemp, err := strconv.ParseInt(nk[0], 10, 64)
		checkError(err)
		n := int32(nTemp)

		kTemp, err := strconv.ParseInt(nk[1], 10, 64)
		checkError(err)
		k := int32(kTemp)

		result := dortmundDilemma(n, k)

		fmt.Fprintf(writer, "%d\n", result)
	}

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
