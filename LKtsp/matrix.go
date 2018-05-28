package main

type Matrix struct {
	d    [][]float64 // distance matrix
	id   []int
	tour []int
}

/*
	vector<int> getCurrentTour();
    double getCurrentTourDistance();
    void optimizeTour();
    void printTour();
	void printTourIds();


	void joinLocations(int i, int j);
    void lkMove(int tourStart);
    void reverse(int start, int end);
    bool isTour();
*/

func (m *Matrix) init() {
	// initialize tour
	size := len(m.id)
	m.tour = make([]int, 0, size)

	// initial 'random' tour
	for i := 0; i < size; i++ {
		m.tour = append(m.tour, (i+1)%size)
	}
}

// result will be in m.tour
func (m *Matrix) solve() {
	var old_distance, new_distance, diff float64
	size := len(m.id)
	for j := 0; j < 100; j++ {
		for i := 0; i < size; i++ {
			m.lkMove(i)
		}
		new_distance = m.getCurrentTourDistance()
		diff = old_distance - new_distance
		if j != 0 {
			//			assert(diff >= 0)
			if diff == 0 {
				//cout << "Converged after " << j << " iterations" << endl;
				break
			}
		}
		old_distance = new_distance
	}
}

func (m *Matrix) lkMove(tourStart int) {
	var broken_set, joined_set [][2]int
	tour_opt := m.tour
	var g_opt, g, g_local, y_opt_length, broken_edge_length, g_opt_local float64

	lastNextV := tourStart
	var nextV, nextFromV, lastPossibleNextV int

	fromV := m.tour[lastNextV]

	var broken_edge [2]int

	//	initialTourDistance := m.getCurrentTourDistance()

	for {
		nextV = -1

		broken_edge = makeSortedPair(lastNextV, fromV)
		broken_edge_length = m.d[broken_edge[0]][broken_edge[1]] //edgeDistances[broken_edge.first][broken_edge.second]

		// Condition 4(c)(1)
		if countPairInSet(joined_set, broken_edge) > 0 {
			break
		}

		for possibleNextV := m.tour[fromV]; nextV == -1 && possibleNextV != tourStart; possibleNextV = m.tour[possibleNextV] {
			// calculate local gain
			g_local = broken_edge_length - m.d[fromV][possibleNextV]

			if !(countPairInSet(broken_set, (makeSortedPair(fromV, possibleNextV))) == 0 &&
				g+g_local > 0 &&
				countPairInSet(joined_set, (makeSortedPair(lastPossibleNextV, possibleNextV))) == 0 &&
				m.tour[possibleNextV] != 0 && // not already joined to start
				possibleNextV != m.tour[fromV]) {
				lastPossibleNextV = possibleNextV
				continue
			}
			nextV = possibleNextV
		}

		if nextV != -1 {
			// add to our broken_set and joined_set
			broken_set = append(broken_set, broken_edge)
			joined_set = append(joined_set, makeSortedPair(fromV, nextV))

			// condition 4(f)
			y_opt_length = m.d[fromV][tourStart] // y_i_opt

			// The tour length if we exchanged the broken edge (x_i)
			// with y_opt, (t_{2i}, t_0)
			g_opt_local = g + (broken_edge_length - y_opt_length)
			if g_opt_local > g_opt {
				g_opt = g_opt_local
				tour_opt = m.tour
				// join the optimal tour
				tour_opt[tourStart] = fromV
			}

			// recalculate g
			g += broken_edge_length - m.d[fromV][nextV]
			// reverse tour direction between newNextV and fromV
			// implicitly breaks x_i
			m.reverse(fromV, lastPossibleNextV)
			nextFromV = lastPossibleNextV
			// build y_i
			m.tour[fromV] = nextV

			// set new fromV to t_{2i+1}
			// and out lastNextV to t_{2i}
			lastNextV = nextV
			fromV = nextFromV

		}
	}

	// join up
	m.tour = tour_opt
	//		distanceAfter := m.getCurrentTourDistance()
	//		assert(distanceAfter <= initialTourDistance)
}

func (m *Matrix) getCurrentTourDistance() float64 {
	return 0.0
}

func makeSortedPair(i int, j int) (pair [2]int) {
	pair[0] = i
	pair[1] = j
	return pair
}

func countPairInSet([][2]int, [2]int) int {
	return 0
}

func (m *Matrix) reverse(start, end int) {
	current := start
	next := m.tour[start]
	var nextNext int
	for {
		nextNext = m.tour[next]

		// reverse the direction at this point
		m.tour[next] = current

		// move to the next pointer
		current = next
		next = nextNext
		if current == end {
			break
		}
	}
}

func (m *Matrix) printTourIds() {
	/*	int current = 0;
		do {
		  cout << ids[current] << endl;
		  current = tour[current];
		} while (current != 0);
	*/
}
