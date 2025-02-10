package main

/*
#cgo CFLAGS: -Iccode
#cgo LDFLAGS: -L. -lmycodelib
#include "ccode/a.h"
*/
import "C"
import "fmt"

func main() {
	C.function_a()
	fmt.Println("!!!")
}
