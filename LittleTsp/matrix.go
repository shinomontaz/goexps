package main

import "fmt"

type Matrix struct {
	arr [][]float64
}

func NewEmptyMatrix(n int) *Matrix {
	matrix := &Matrix{}
	matrix.arr = make([][]float64, 0, n)

	for i := 0; i < n; i++ {
		matrix.arr[i] = make([]float64, 0, n)
	}

	return matrix
}

func NewMatrix(arr [][]float64) *Matrix {
	matrix := &Matrix{}

	// do a copy
	matrix.arr = make([][]float64, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		matrix.arr[i] = append(matrix.arr[i], arr[i]...)
	}
	return matrix
}

func (m *Matrix) GetTransposed() *Matrix {
	newMatrix := NewMatrix(m.arr)

	for i, row := range m.arr {
		for j, el := range row {
			newMatrix.arr[j][i] = el
		}
	}
	return newMatrix
}

func (m *Matrix) GetDimension() int {
	return len(m.arr)
}

// Adduct function:
// 1 adduction by rows
// 2 adduction by cols
// preverse weight of zeroes
func (m *Matrix) Adduct() (adducted *Matrix, est_i []float64, est_j []float64) {
	adducted = NewMatrix(m.arr)

	for i, row := range m.arr {
		// find the minimum of row
		rowMin := SliceMin(row)
		est_i[i] = rowMin
		for j, el := range row {
			adducted.arr[i][j] = el - rowMin
		}
	}

	transponed := adducted.GetTransposed()

	for j, col := range transponed.arr {
		colMin := SliceMin(col)
		est_j[j] = colMin
		for i := range col {
			adducted.arr[i][j] -= colMin
		}
	}

	return adducted, est_i, est_j
}

func (m *Matrix) GetZeroWeight( [2]int ) {
    if( !isset( $this->arr[$i] ) || !isset( $this->arr[$i][$j] ) || $this->arr[$i][$j] != 0 ) {
      throw new Exception( 'getZeroWeight receives an element, that is not zero in matrix or not exists' );
    }
    // get min element in row $zero[0]
    // get min element in col $zero[1]
    // sum them
    $row_min = self::_INF;
    $col_min = self::_INF;

    foreach( $this->arr[$i] as $index => $elem ) {
      if( $index == $j ) continue;
      $row_min = ( $row_min > $elem ) ? $elem : $row_min;
    }

    foreach( $this->arr as $index => $row ) {
      if( $index == $i ) continue;
      $elem = $row[$j];
      $col_min = ( $col_min > $elem ) ? $elem : $col_min;
    }

    $weight = $col_min + $row_min;

    return $weight;
  }

func (m *Matrix) GetZeroes( ) [][2]int {
    zeroes = make([][2]int)
	for i, row := range m.arr {
      for j, el := range row {
        if el == 0.0
          zeroes = append( zeroes, [2]int{ i, j } );
      }
    }
    
    return zeroes
  }

  func GetZeroWeight

func (m *Matrix) GetHeaviest() [2]int {
	zeroes := m.GetZeroes();
		
	max := 0.0
	heaviest, exists := zeroes[0]
	if !exists {
		return [2]int{}
	}
	for zero := range zeroes {
		candidate_weight := m.GetZeroWeight( zero )
		if( max < candidate_weight ) {
			max = candidate_weight
			heaviest = zero
		  }
		}
	
		return heaviest
}

func (m *Matrix) GetPoints() [][2]int {
	points = make([][2]int, len(m.arr))
	for i, row := range m.arr {
		for j, el := range row {
			points = append(points, [2]int{i, j})
		}
	}
	return points
}

func (m *Matrix) GetElement(i, j) (float64, error) {
	if row, exists := m.arr[i]; exists {
		if el, exists2 := row[j]; exists2 {
			return el, nil
		}
	}
	return nil, fmt.Errorf("no element")
}

func SliceMin(in []float64) float64 {
	min := in[0]
	for _, v := range in {
		if min > v {
			min = v
		}
	}
	return min
}

func SliceMax(in []float64) float64 {
	max := in[0]
	for _, v := range in {
		if max < v {
			max = v
		}
	}
	return max
}
