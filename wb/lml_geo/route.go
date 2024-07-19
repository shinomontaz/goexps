package main

func fitness(routes [][]int, tm [][]float64) (float64, float64) {
	f := 0.0
	o := 0.0
	for _, r := range routes {
		ff, ov := cost(r, tm)
		f += ff
		o += ov
	}
	return f, o
}

func max_overtime(routes [][]int, tm [][]float64) (float64, float64, int) {
	max := 0.0
	mean := 0.0
	num := 0
	for _, r := range routes {
		_, ov := cost(r, tm)
		if ov > max {
			max = ov
		}
		if ov > 0 {
			mean += ov
			num++
		}
	}
	mean = mean / float64(num)
	return max, mean, num
}

func cost(route []int, tm [][]float64) (float64, float64) {
	time := 0.0
	overtime := 0.0
	for i := 1; i < len(route); i++ {
		time += tm[route[i-1]][route[i]]
		if tm[route[i-1]][route[i]] > 0 {
			time += box_time
		}
	}

	if len(route) > 1 {
		// добаляем время на возврат на склад
		time += tm[route[len(route)-1]][0]
	}

	if time > courier_shift {
		overtime = time - courier_shift
	}

	return time + courier_shift, overtime
}
