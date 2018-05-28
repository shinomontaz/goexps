package solver

import (
	"fmt"
	"testing"
)

func BenchmarkGsolver_Solve(b *testing.B) {
	points := createPoints(10)
	distances := calcDistances(points)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := New(points, distances)
		_ = s.Solve(0.01, 100, 100000)
	}
}

func TestSmoke(t *testing.T) {
	points := createPoints(10)
	distances := calcDistances(points)
	s := New(points, distances)
	res := s.Solve(0.01, 100, 100000)

	fmt.Println(res)
}
