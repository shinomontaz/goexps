package main

/*
#cgo CFLAGS: -Iccode
#cgo LDFLAGS: -Lccode/OBJ -llkh
#include "INCLUDE/LKH.h"
*/
import "C"
import "fmt"

func main() {
	fmt.Println(C.Random())
}
