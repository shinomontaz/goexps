package main

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Number     string
	Items      []int
	Created_at time.Time
}

var GoodIds = []int{1, 2, 3, 4, 5}

func newOrder() Order {
	n := rnd.Intn(len(GoodIds))
	perm := rnd.Perm(len(GoodIds))

	o := Order{
		Number:     uuid.New().String(),
		Items:      make([]int, n),
		Created_at: time.Now(),
	}

	for i := range n {
		o.Items[i] = GoodIds[perm[i]]
	}

	return o
}
