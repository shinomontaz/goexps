package tsp

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkHeaps(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	// Define the size of the distance matrix
	size := 10
	first := 1
	last := 1

	// Generate the distance matrix
	dm := generateDistanceMatrix(size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := solve(dm, first, last)
		if err != nil {
			b.Fatalf("Exec3 failed: %v", err)
		}
	}
}

func generateDistanceMatrix(size int) [][]int64 {
	dm := make([][]int64, size)
	for i := range dm {
		dm[i] = make([]int64, size)
		for j := range dm[i] {
			if i != j {
				dm[i][j] = int64(rand.Intn(100) + 1) // Random distance between 1 and 100
			}
		}
	}
	return dm
}

// func TestExec2(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		dm      [][]int64
// 		first   int
// 		last    int
// 		want    []int32
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid case with 4 points",
// 			// dm: [][]int64{
// 			// 	{0, 1, 2, 3},
// 			// 	{1, 0, 4, 5},
// 			// 	{2, 4, 0, 6},
// 			// 	{3, 5, 6, 0},
// 			// },
// 			dm: [][]int64{
// 				{0, 10, 15},
// 				{10, 0, 20},
// 				{15, 20, 0},
// 			},
// 			first:   0,
// 			last:    0,
// 			want:    []int32{0, 1, 2, 0},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &algo{}
// 			got, err := a.Exec2(tt.dm, tt.first, tt.last)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Exec2() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Exec2() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
