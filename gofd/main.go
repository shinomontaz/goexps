package main

import (
	"fmt"

	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

func main() {
	store = core.CreateStore()
	n := 8

	queens := make([]core.VarId, n)
	for i := 0; i < n; i++ {
		queens[i] = core.CreateIntVarFromTo(fmt.Sprintf("Queen - %d", i), store, 0, n-1)
	}

	prop := propagator.CreateAlldifferent(queens...)
	store.AddPropagators(prop)

	// Диагональные ограничения
	left_offset := make([]int, len(queens))
	right_offset := make([]int, len(queens))
	for i, _ := range queens {
		left_offset[i] = -i
		right_offset[i] = i
	}

	left_prop := interval.CreateAlldifferentOffset(queens, left_offset) // тоже самое, что и AllDifferent, но учитывает доп. переменную left_offset
	/*
		queens[i] + left_offset[i] != queens[j] + left_offset[j] для любого i != j

		0 + 0 != 1 - 1 // i = 0, j = 1 queens[0] = 0, queens[1] = 1
		0 + 0 != 3 - 2 // i = 0, j = 2 queens[0] = 0, queens[2] = 3

	*/
	store.AddPropagator(left_prop)
	right_prop := interval.CreateAlldifferentOffset(queens, right_offset)
	store.AddPropagator(right_prop)

	query := labeling.CreateSearchAllQuery()
	solutionFound := labeling.Labeling(store, query, labeling.SmallestDomainFirst, labeling.InDomainMin)
	if solutionFound {
		fmt.Println(n, "queens problem has", len(query.GetResultSet()), "solutions.")
	}

}
