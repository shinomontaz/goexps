package main

//#cgo CFLAGS: -g
//#include <stdlib.h>
import "C"
import "fmt"

func Random() int {
	return int(C.rand())
}

func main() {
	fmt.Println(Random())
}
