package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var str []byte
var states [][]byte

func main() {
	//Enter your code here. Read input from STDIN. Print output to STDOUT
	str = make([]byte, 0)

	var cmd string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cmd = scanner.Text()

	ops_num, _ := strconv.Atoi(cmd)

	states = make([][]byte, 0)

	states = append(states, []byte{})

	for i := 0; i < ops_num; i++ { // endless loop to get input data
		scanner.Scan()
		cmd = scanner.Text()
		process(cmd)
	}
}

func process(cmd string) {
	args := strings.Split(cmd, " ")
	op, _ := strconv.Atoi(args[0])

	switch op {
	case 1: // append
		apnd([]byte(args[1]))
		curr_str := make([]byte, len(str))
		for i, b := range str {
			curr_str[i] = b
		}
		states = append(states, curr_str)
	case 2: // delete
		idx, _ := strconv.Atoi(args[1])
		del(idx)
		curr_str := make([]byte, len(str))
		for i, b := range str {
			curr_str[i] = b
		}
		states = append(states, curr_str)
	case 3: // print
		idx, _ := strconv.Atoi(args[1])
		prnt(idx)
	case 4: // undo
		states = states[:len(states)-1]
		if len(states) > 0 {
			str = states[len(states)-1]
		}
	}
}

func apnd(w []byte) {
	str = append(str, w...)
}

func del(i int) {
	str = str[:len(str)-i]
}

func prnt(i int) {
	fmt.Printf("%c\n", str[i-1])
}

func und(str *[]byte) {
	// undo last command
}
