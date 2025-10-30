package bnb

import (
	"math"
)

type Point struct {
	weight  int64 // стоимость этой точки в текущей мартице
	penalty int64 // штрафы нулевой точки
	row     int   // абс. позиция в локальной матрице
	col     int   //
	ir      int   // локальные индексы текущего уровня
	ic      int
}

const (
	MaxDimension = 96
	sizeCube     = MaxDimension * MaxDimension * MaxDimension
	sizeSquare   = MaxDimension * MaxDimension
)

// branches -- поиск оптимального пути методом "ветви и границы"
// @return слайс порядка обхода и суммарное расстояние
// реализация без рекурсии и с предварительным выделением памяти, вложенным методами для работы со стеком
func branches(distances [][]int64, start, end int) ([]int, int64) {
	lenD := len(distances)

	if lenD < 2 {
		return nil, 0
	}
	if lenD > MaxDimension {
		return nil, 0.0 // Передельная матрица на стеке 96х96х96х8=6.75Мб.
	}

	// 0. формируем наборы данных:
	var (
		matrix     [sizeCube]int64
		totals     [MaxDimension]int64 // общая стоимость текущего уровня с выбранным отрезком
		rowMins    [sizeSquare]int64   // текущие миниумы строк и столбцов
		colMins    [sizeSquare]int64
		rowIds     [sizeSquare]int
		colIds     [sizeSquare]int
		levels     [MaxDimension]int   // начальные смещения уровней в матрице
		ofsets     [MaxDimension]int   // смещения в линейных  массивах по уровням
		sizes      [MaxDimension]int   // размеры массивов по уровням
		selecteds  [MaxDimension]Point // выбрано: вес, штраф + обе позиции
		bestPoints [MaxDimension]Point
		bestCosts  [MaxDimension]int64 // оценки лучшего пути по уровням
		rootCost   = int64(0)
		lastPoints [2]Point // последние 2 точки отдельно из 2х2

		// функция расчета миниумов по строкам и стобцам сразу с вычитанием миниума по строкам
		calcMins = func(level int) {
			ofM := levels[level]
			ofL := ofsets[level]
			size := sizes[level]

			// поиск миниума по строкам
			next := 0
			pos := 0
			min := int64(math.MaxInt64)
			for of := ofM; of < ofM+size*size; of++ {
				mRow := matrix[of]
				if mRow < int64(math.MaxInt64) {
					if mRow < min {
						min = mRow
					}
				}
				next++
				if next == size {
					rowMins[ofL+pos] = min
					next = 0
					pos++
					min = int64(math.MaxInt64)
				}
			}
			// редукция строк и поиск миниума по столбцам
			next = 0
			pos = 0
			min = int64(math.MaxInt64)
			for of := ofM; of < ofM+size*size; {
				mCol := matrix[of]
				if mCol < int64(math.MaxInt64) {
					mCol -= rowMins[ofL+next]
					if mCol < min {
						min = mCol
					}
				}
				of += size
				next++
				if next == size {
					colMins[ofL+pos] = min
					pos++
					if pos == size {
						break
					}
					of = ofM + pos
					next = 0
					min = int64(math.MaxInt64)
				}
			}
		}

		// функция вычисления нижней границы
		calcCost = func(level int) int64 {
			cost := int64(0)

			of := ofsets[level]
			for i := of; i < of+sizes[level]; i++ {
				rowMin := rowMins[i]
				colMin := colMins[i]
				if rowMin < int64(math.MaxInt64) { // исключенная строка не участвует в стоимости
					cost += rowMin
				}
				if colMin < int64(math.MaxInt64) { // исключенный столбец не участвует в стоимости
					cost += colMin
				}
			}

			return cost
		}

		// функция вычисления штрафов нулевых клеток, поиск нулевой точки с максимальным штрафом
		findZeroMax = func(level int) {
			ofM := levels[level]
			ofL := ofsets[level]
			size := sizes[level]
			maxCost := int64(-1)

			for r := 0; r < size; r++ {
				rowPenalty := rowMins[ofL+r]
				for c := 0; c < size; c++ {
					of := ofM + r*size + c
					cur := matrix[of]
					if cur < int64(math.MaxInt64) { // не оцениваем точки с мах!
						colPenalty := colMins[ofL+c]
						if cur-rowPenalty-colPenalty == 0 {
							// ищем локальные миниумы строки и столбца
							colMin := int64(math.MaxInt64)
							rowMin := int64(math.MaxInt64)
							for rc := 0; rc < size; rc++ {
								ofr := ofM + r*size + rc
								ofc := ofM + rc*size + c
								if rowMin > 0 && c != rc { // меньше 0 не бывает!
									mRow := matrix[ofr]
									if mRow < int64(math.MaxInt64) {
										mRow = mRow - (rowPenalty + colMins[ofL+rc])
										if mRow < rowMin {
											rowMin = mRow
										}
									}
								}
								if colMin > 0 && r != rc {
									mCol := matrix[ofc]
									if mCol < int64(math.MaxInt64) {
										rpi := ofL + rc
										rp := rowMins[rpi]
										mCol2 := mCol - (rp + colPenalty) // тоже, строка - текущая!
										if mCol2 < 0 {
											panic("mCol < 0")
										}
										if mCol2 < colMin {
											colMin = mCol2
										}
									}
								}
								if rowMin == 0 && colMin == 0 { // оба нули, уходим
									break
								}
							}

							cost := rowMin + colMin

							if maxCost < cost { // первая из равных!
								maxCost = cost
								selecteds[level] = Point{matrix[of], cost, rowIds[ofL+r], colIds[ofL+c], r, c}
							}
						}
					}
				}
			}
		}

		// запрещение куска пути, обратного выбранной точке
		// @TODO: как-то бы ускорить поиск цепочек..
		closeBackWay = func(level, origRow, origCol int) {
			// 1. ищем обратный путь среди выбранного
			for i := level - 2; i >= 0; i-- {
				if selecteds[i].row == origCol {
					origCol = selecteds[i].col
					i = level - 1 // повтор сначала! Цепочка м.б. ниже!
				}
			}
			for i := level - 2; i >= 0; i-- {
				if selecteds[i].col == origRow {
					origRow = selecteds[i].row
					i = level - 1
				}
			}

			// 2. Ищем транспонированную точку к заданной
			backRow := -1
			backCol := -1

			ofM := levels[level]
			ofL := ofsets[level]
			size := sizes[level]
			for rc := ofL; rc < ofL+size; rc++ {
				if rowIds[rc] == origCol { // ищем транспонированную точку!
					backRow = rc - ofL
				}
				if colIds[rc] == origRow {
					backCol = rc - ofL
				}
				if backRow >= 0 && backCol >= 0 {
					break
				}
			}

			// если нашли, то запрещаем обратный путь
			if backRow >= 0 && backCol >= 0 {
				of := ofM + backRow*size + backCol
				matrix[of] = int64(math.MaxInt64)
			} else {
				panic("обратный путь не найден") // сюда не должны попадать вовсе.
			}
		}

		// запрещение selected точки в матрице
		closeSelected = func(level int) {
			point := selecteds[level]
			matrix[levels[level]+sizes[level]*point.ir+point.ic] = int64(math.MaxInt64)
		}

		// Редукция матрицы - переход на след. уровень итераций с уменьшением размерности
		reduceMatrix = func(level int) {
			oldOfM := levels[level]
			oldSize := sizes[level]
			oldOfL := ofsets[level]

			subRow := selecteds[level].ir
			subCol := selecteds[level].ic

			level++
			newOfM := levels[level]
			newSize := sizes[level]
			newOfL := ofsets[level]

			// 1. копируем матрицу с коррекцией стоимостей и индексы оригинальных точек без строки и столбца:
			nr := 0
			for r := 0; r < oldSize; r++ {
				if r == subRow {
					continue
				}

				nc := 0
				for c := 0; c < oldSize; c++ {
					if c == subCol {
						continue
					}
					oldOf := oldOfM + r*oldSize + c
					newOf := newOfM + nr*newSize + nc

					curM := matrix[oldOf]
					if curM != int64(math.MaxInt64) {
						curM -= rowMins[oldOfL+r] + colMins[oldOfL+c]
					}
					matrix[newOf] = curM
					colIds[newOfL+nc] = colIds[oldOfL+c]
					nc++
				}
				rowIds[newOfL+nr] = rowIds[oldOfL+r]
				nr++
			}

			// 2. запрещаем обратный путь от нулевой точки по оригинальной позиции!
			closeBackWay(level, rowIds[oldOfL+subRow], colIds[oldOfL+subCol])

			// 4. Переносим стоимость отрезков
			totals[level] = totals[level-1]
		}

		// остаток 2х2, без возврата 1 путь "не М" (выбранная точка), с возвратом "диагональ выбранной точки"
		getLast = func(level int) {
			ofM := levels[level]
			ofL := ofsets[level]
			size := sizes[level]

			point := selecteds[level]
			lastPoints[0] = point

			if start == end {
				row1 := 1 - point.ir
				col1 := 1 - point.ic

				lastPoints[1] = Point{
					matrix[ofM+row1*size+col1], -1, rowIds[ofL+row1], colIds[ofL+col1], row1, col1,
				}
			}
		}

		// перенос найденного пути как лучшего
		storeBest = func() {
			var i int
			for i = 0; i < lenD-2; i++ {
				bestPoints[i] = selecteds[i]
				bestCosts[i] = totals[i]
			}
			// i == lenD-2 !!!
			bestPoints[i] = lastPoints[0]
			bestCosts[i] = totals[i]

			// Нужно след. значение для сравнения оценок далее. Дублируем последнее.
			bestCosts[i+1] = totals[i]

			if start == end {
				i++
				bestPoints[i] = lastPoints[1]
			}
		}

		// преобразование пути в маршрут и вычисление его стоимости
		makeRoute = func(points [MaxDimension]Point, lenD int) ([]int, int64) {
			lenP := lenD - 1
			if start == end {
				lenP = lenD // с возвратом на 1 точку маршрут длинее!
			}
			// @TODO: убрать выделение памяти .. не нужно по идее.
			// маршрут больше числа кусков на 1шт.
			route := make([]int, 0, lenP+1)

			num := 0
			total := int64(0)
			for i := 0; i < lenP; i++ {
				ok := false
				for _, p := range points[:lenP] {
					if p.row == num {
						total += distances[num][p.col]
						route = append(route, num)
						num = p.col
						ok = true

						break
					}
				}
				if !ok {
					panic("разрыв в пути!")
				}
			}
			route = append(route, num) // последнюю точку искать не зачем..

			return route, total
		}
	)

	// 1. Заполняем начальный уровень без точки (0,0)
	ofM := 0
	ofL := 0
	for r := 0; r < lenD; r++ {
		for c := 0; c < lenD; c++ { // без возврата в 0,0 - колонка 0 === М ниже
			of := r*lenD + c
			if r == c {
				matrix[of] = int64(math.MaxInt64)
			} else {
				if distances[r][c] <= 0.0 {
					panic("d <=0 ")
				}
				matrix[of] = distances[r][c]
			}
		}

		if start != end {
			matrix[r*lenD] = int64(math.MaxInt64) // запрещаем возврат в начало (0,0)
		}

		rowIds[r] = r
		colIds[r] = r

		levels[r] = ofM
		ofsets[r] = ofL
		sizes[r] = lenD - r

		ofM += (lenD - r) * (lenD - r)
		ofL += lenD - r
	}

	// 2. Первый проход, спускаемся без сокращения матрицы ("не включая"):

	level := 0
	cost := int64(0)
	for level = 0; level < lenD-1; level++ {
		// расчет уровня ВКЛЮЧАЯ точку:
		calcMins(level)
		cost = calcCost(level)
		totals[level] += cost
		if cost == int64(math.MaxInt64) {
			totals[level] = int64(math.MaxInt64)
		}
		if level == 0 {
			// @TODO если без возврата, то надо исключить одно значение из rowMins .. какое?!?
			rootCost = cost
		}
		findZeroMax(level)

		// редуцируем матрицу, без последней итерации 2х2 уменьшать незачем!
		if level < lenD-2 {
			reduceMatrix(level)
		} else {
			// остаток 2х2 отдельно, тут пара точек с макс штрафом
			getLast(level)
		}
	} // level = lenD-1!

	// 3. Собираем найденный путь и его стоимость как образец
	storeBest()

	if bestCosts[level-1] <= rootCost {
		// это один из оптимальных путей, лучше не будет
		return makeRoute(bestPoints, lenD)
	}

	// 4. Ищем путь лучше подъемом по графу снизу вверх. Обнуляем стоимость и запрещаем выбранную точку уровня -2.
	num := 0
	level -= 2
	totals[level] = totals[level-1]
	closeSelected(level)
	for {
		num++
		calcMins(level)

		total := totals[level]
		total += calcCost(level)

		if total >= bestCosts[level+1] { // не лучше чем ниже на уровень, редуцировать незачем.
			if level > 0 {
				level--
				if level <= 0 {
					totals[0] = 0 // пересчет заново, самый верх
				} else {
					totals[level] = totals[level-1]
				}
				// поднялись на уровень и закрыли точку с учетом транспонирования внутри
				closeSelected(level)
			} else {
				break // это самый верх, лучше уже не будет
			}
		} else {
			totals[level] = total
			findZeroMax(level)

			if level < lenD-2 {
				reduceMatrix(level)
				level++ // это не конец, ищем редуцированием матрицы
			} else {
				getLast(level)
				storeBest()
				if bestCosts[level+1] <= rootCost {
					break
				}
				level-- // не самый лучший, поднимаемся выше
			}
		}
	}

	return makeRoute(bestPoints, lenD)
}
