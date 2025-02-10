package main

/*
#cgo CFLAGS: -Iccode
#cgo LDFLAGS: -L. -lmy
#include <stdlib.h>
#include <stdio.h>
#include "LKH.h"
#include "my.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	params := C.createDefaultMyParameters()

	params.ExtraCandidateSetType = C.POPMUSIC
	params.InitialTourAlgorithm = C.CTSP_ALG
	params.MTSPObjective = C.MINSUM
	params.Recombination = C.GPX2
	params.SubproblemSpecial = C.SubproblemSpecialEnum(C.SUBPROBLEM_KARP)
	params.SubproblemSpecial2 = C.SubproblemSpecial2Enum(C.SPECIALSUBPROBLEM_COMPRESSED)

	fmt.Printf("ExtraCandidateSetType: %d\n", params.ExtraCandidateSetType)
	fmt.Printf("SubproblemSpecial2: %d\n", params.SubproblemSpecial2)

	var problem C.MyProblem
	problem.Capacity = 100
	problem.DemandDimension = 1
	problem.MTSPDepot = 1
	problem.ProblemType = C.TSP
	problem.Dimension = 10
	problem.EdgeWeightType = C.EUC_2D

	problem.nodeCoords = (*C.NodeCoord)(C.malloc(C.size_t(10) * C.size_t(unsafe.Sizeof(C.NodeCoord{}))))
	defer C.free(unsafe.Pointer(problem.nodeCoords))

	for i := 0; i < 10; i++ {
		(*C.NodeCoord)(unsafe.Pointer(uintptr(unsafe.Pointer(problem.nodeCoords)) + uintptr(i)*unsafe.Sizeof(C.NodeCoord{}))).Id = C.int(i + 1)
		(*C.NodeCoord)(unsafe.Pointer(uintptr(unsafe.Pointer(problem.nodeCoords)) + uintptr(i)*unsafe.Sizeof(C.NodeCoord{}))).X = C.double(float64(i) + 1.0)
		(*C.NodeCoord)(unsafe.Pointer(uintptr(unsafe.Pointer(problem.nodeCoords)) + uintptr(i)*unsafe.Sizeof(C.NodeCoord{}))).Y = C.double(float64(i) + 2.0)
		(*C.NodeCoord)(unsafe.Pointer(uintptr(unsafe.Pointer(problem.nodeCoords)) + uintptr(i)*unsafe.Sizeof(C.NodeCoord{}))).Z = C.double(float64(i) + 3.0)
	}

	fmt.Println(problem)
}
