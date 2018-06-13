package main

import (
	"fmt"
	"testing"
)

func TestDistance(t *testing.T) {

	p := Point{x: 0, y: 0}
	q := Point{x: 1, y: 0}

	fmt.Println(distance(p, q))
}
