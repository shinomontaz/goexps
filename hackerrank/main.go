package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	n := int32(nTemp)

	var puzzle [][]int32
	for i := 0; i < int(n); i++ {
		puzzleRowTemp := strings.Split(readLine(reader), " ")

		var puzzleRow []int32
		for _, puzzleRowItem := range puzzleRowTemp {
			puzzleItemTemp, err := strconv.ParseInt(puzzleRowItem, 10, 64)
			checkError(err)
			puzzleItem := int32(puzzleItemTemp)
			puzzleRow = append(puzzleRow, puzzleItem)
		}

		if len(puzzleRow) != int(int(n)) {
			panic("Bad input")
		}

		puzzle = append(puzzle, puzzleRow)
	}

	// Write Your Code Here
	rand.Seed(time.Now().UnixNano())
	// movesCounter := rand.Intn(500)
	// moves := make([][]int, movesCounter)
	// for i := 0; i < movesCounter; i++ {
	//     randI := rand.Intn(int(n))
	//     randJ := rand.Intn(int(n))
	//     randK := rand.Intn(int(n))
	//     moves = append(moves, []int{randI, randJ, randK})
	// }

	// prepare dump output

	initSol := &Sol{
		puzzle: puzzle,
		moves:  make([][]int, 0),
	}

	finalSol := sa(initSol)

	finalSol.Print()
}

func move(puzzle [][]int32, i, j, k int) [][]int32 {
	subpuzzle := make([][]int32, k)
	for idx := 0; idx < k; idx++ {
		subpuzzle[idx] = make([]int32, k)
		for jdx := 0; jdx < k; jdx++ {
			subpuzzle[idx][jdx] = puzzle[i+idx][j+jdx]
		}
	}

	// make clockwise subpuzzle move
	rotated := make([][]int32, k)
	for idx := 0; idx < k; idx++ {
		rotated[idx] = make([]int32, k)
		for jdx := 0; jdx < k; jdx++ {
			rotated[idx][jdx] = subpuzzle[k-jdx-1][idx]
		}
	}

	// fmt.Println("subpuzzle", subpuzzle, i, j, k)
	// fmt.Println("rotated", rotated)

	// append rotated subpuzzle to puzzle in position
	for idx := i; idx < k; idx++ {
		//		fmt.Println("idx, i", idx, i)
		for jdx := j; jdx < k; jdx++ {
			//			fmt.Println("jdx, j", jdx, j)
			puzzle[idx][jdx] = rotated[idx-i][jdx-j]
		}
	}

	return puzzle
}

func goodness(puzzle [][]int32) int {
	return simplegood(puzzle) + simplegood(transpone(puzzle))
}

func transpone(m [][]int32) [][]int32 {
	t := make([][]int32, len(m))
	for i := range m { // by rows
		t[i] = make([]int32, len(m))
	}
	for i, row := range m { // by rows
		for j, val := range row {
			t[j][i] = val
		}
	}
	return t
}

func simplegood(puzzle [][]int32) (score int) {
	// find all good pairs by rows and cols

	for _, row := range puzzle { // by rows
		for _, i := range row {
			for ii := i; ii < int32(len(row)); ii++ {
				if ii > i {
					score++
				}
			}
		}
	}

	return score
}

type Sol struct {
	puzzle [][]int32
	moves  [][]int
}

func (s *Sol) Print() {
	fmt.Println(len(s.moves))
	for _, move := range s.moves {
		fmt.Printf("%d %d %d\n", move[0]+1, move[1]+1, move[2])
	}
}

func (s *Sol) Fitness() float64 {
	return 1 / float64(goodness(s.puzzle))
}

func (s *Sol) Copy() *Sol {
	newSol := Sol{
		puzzle: make([][]int32, len(s.puzzle)),
		moves:  make([][]int, len(s.moves)),
	}

	for i, row := range s.puzzle {
		newSol.puzzle[i] = make([]int32, len(row))
		for j, val := range row {
			newSol.puzzle[i][j] = val
		}
	}

	for i, move := range s.moves {
		newSol.moves[i] = make([]int, len(move))
		for j, val := range move {
			newSol.moves[i][j] = val
		}
	}

	return &newSol
}

func (s *Sol) Mutate() *Sol {
	newSol := s.Copy()
	randI, randJ := rand.Intn(len(s.puzzle)), rand.Intn(len(s.puzzle))
	randK := 1 + rand.Intn(len(s.puzzle)-int(math.Max(float64(randI), float64(randJ))))
	newSol.puzzle = move(newSol.puzzle, randI, randJ, randK)
	newSol.moves = append(newSol.moves, []int{randI, randJ, randK})

	return newSol
}

func sa(currSol *Sol) *Sol {
	T := 1.0
	Tmin := 0.001
	cooling := 0.999
	oldEnergy := currSol.Fitness()
	steps := 0
	maxSteps := 500

	for T > Tmin && steps <= maxSteps {
		newSol := currSol.Mutate()
		newEnergy := newSol.Fitness()
		if newEnergy < oldEnergy {
			currSol = newSol
			oldEnergy = newEnergy
			steps++
		} else {
			dice := rand.Float64()
			if dice > saAcceptance(T, oldEnergy, newEnergy) {
				currSol = newSol
				oldEnergy = newEnergy
				steps++
			}
		}
		T *= cooling
	}

	return currSol
}

func saAcceptance(T float64, oldEnergy, newEnergy float64) float64 {
	return math.Exp((newEnergy - oldEnergy) / T)
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
