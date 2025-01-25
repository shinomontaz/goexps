package main

/*
#cgo CFLAGS: -I./ccode
#cgo LDFLAGS: -L./ccode -ladd
#include "add.h"
*/
import "C"
import "fmt"

func main() {
	a := 6
	b := 3
	result := C.add(C.int(a), C.int(b))
	fmt.Printf("The sum of %d and %d is %d\n", a, b, result)
}
