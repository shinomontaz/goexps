package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {

	n := 10
	points := createPoints(n)
	dMatrix := calcDistances(points)

	route := make([][2]int, 0, n)

	mainMatrix := NewMatrix(dMatrix)

	modMatrix, est_i, est_j := mainMatrix.Adduct()
	recursive(modMatrix, route)

	// reorder route to have it consistent with beginnnung in $start

	subroute := make([][2]int, 0, n)
	subroute = append(subroute, route[0])

	route2 := make([][2]int, 0, n)
	for i, point := range route {
		route2 = append(route2, point)
	}

	buildSubRouteFromPoint(route[0], subroute, route2)

	for subroute[0][0] != 0 {
		shiftRouteLeft(subroute)
	}

	route = subroute

	fmt.Println(route)
}

func shiftRouteLeft(route [][2]int) {
	route = append(route[1:], route[0])
}

func buildSubRouteFromPoint(point [2]int, subroute [][2]int, route [][2]int) { // recursive function
	pair, err := findNextPoint(point, route)
	if err == nil {
		subroute = append(subroute, pair)
		for _, subpoint := range subroute {
			index := findIndexOfPoint(subpoint, route) // check the error!

			route = route[:index+copy(route[index:], route[index+1:])]
		}
		buildSubRouteFromPoint(pair, subroute, route)
	}
}

func findIndexOfPoint(point [2]int, route [][2]int) int {
	for i, pp := range route {
		if pp[0] == point[0] && pp[1] == point[1] {
			return i
		}
	}
	return -1
}

func findNextPoint(point [2]int, route [][2]int) ([2]int, error) { // find a point in route that connected to one
	for i, pp := range route {
		if point[1] == pp[0] {
			return pp, nil
		}
	}
	return [2]int{}, fmt.Errorf("no result")
}

type LatLng struct {
	Lat float64
	Lng float64
}

func createPoints(n int) []*LatLng {
	res := make([]*LatLng, 0)
	for i := 0; i < n; i++ {
		res = append(res, &LatLng{
			Lat: rand.Float64() * 100,
			Lng: rand.Float64() * 100,
		})
	}
	return res
}

func calcDistances(points []*LatLng) [][]float64 {
	res := make([][]float64, 0)
	for _, from := range points {
		row := make([]float64, 0)
		for _, to := range points {
			row = append(row, getDistance(from, to))
		}
		res = append(res, row)
	}

	return res
}

func getDistance(from, to *LatLng) float64 {
	if from == to {
		return 0
	}
	return math.Sqrt(math.Pow(from.Lat-to.Lat, 2) + math.Pow(from.Lng-to.Lng, 2))
}

func recursive( matrix *Matrix ) [][2]int {
  if( matrix.GetDimension() == 1 ) {
	rows := GetKeys( matrix.arr )
  	cols := array_keys( matrix.arr[0] );
  	return [][2]int{ [2]int{rows[0], cols[0]} };
	}

// if matrix consists of 1 element only => finish it
	
	linesToDelete := pointsToEliminate(route, matrix.GetPoints())

	for pointToDelete := range linesToDelete {
		if _, err := matrix.GetElement( pointToDelete[0],  pointToDelete[1] ); err == nil {
			matrix.SetElement( pointToDelete[0], pointToDelete[1], math.MaxFloat64 );
		}
		if _, err := matrix.GetElement( pointToDelete[1], pointToDelete[0] ); err == nil {
			matrix.SetElement( pointToDelete[1], pointToDelete[0], math.MaxFloat64 )
		}
	}
	
	heaviest := matrix.GetHeaviest( ); // detect heavies zero
// change $heaviest to INF
if( $matrix->getElement( $heaviest[0], $heaviest[1] ) !== FALSE ) {
  $matrix->setElement( $heaviest[0], $heaviest[1], DistanceMatrix::_INF );
}

// 1. exclude an element from route: - estimate full matrix and get its estimation. At this step output matrix will be equal to input, but with the element setted to infinity
// 2. include element in route - estimate reduced matrix and get its estimation
$estimation_full_i = array();
$estimation_full_j = array();

$this->adduction( $matrix, $estimation_full_i, $estimation_full_j ); // we have no need to store a adducted matrix
$estimation_full = array_sum($estimation_full_j) + array_sum( $estimation_full_i);

$matrixReduced = DistanceMatrix::init( $matrix->getArray() );

$matrixReduced->unsetRowAndColumn( $heaviest[0], $heaviest[1] );
if( $matrixReduced->getElement( $heaviest[1], $heaviest[0] ) !== FALSE ) {
  $matrixReduced->setElement($heaviest[1], $heaviest[0], DistanceMatrix::_INF );
}

$estimation_reduced_i = array();
$estimation_reduced_j = array();

$matrixReduced2 = $this->adduction( $matrixReduced, $estimation_reduced_i, $estimation_reduced_j );

$estimation_reduced = array_sum($estimation_reduced_i) + array_sum( $estimation_reduced_j );

$endMatrix = '';
if( $estimation_reduced > $estimation_full ) {  // here we choose a full, cause it is min
  $endMatrix = DistanceMatrix::init( $matrix->getArray() ); // element is not in route, so exclude it from further calculations
}
else {
   $route[] = $heaviest;
   $endMatrix = DistanceMatrix::init( $matrixReduced2->getArray() );
}

if( $endMatrix->getDimension() > 2 ) { // remove possible routes that can result in non-hamilton subcycle
}

$this->recursive( $endMatrix, $route );
}

func GetKeys( [][]float64 ) []int { // func GetKeys( [][float64] ) []int 
	mymap := make(map[int]string)
    keys := make([]int, 0, len(mymap))
    for k := range mymap {
        keys = append(keys, k)
	}
	return keys
}

func pointsToEliminate( array $route, array $points ) { // return an array of points to eliminate from further countings
	//
	if( !$route ) {
		return [];
	}
	$eliminate = array();
	foreach( $points as $point ) {
		$testRoute = array_merge( $route, array( $point ) );
		if( $this->isCycleRoute($testRoute) ) {
			$eliminate[] = $point;
		}
	}
	return $eliminate;
}