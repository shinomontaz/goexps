package main

import "fmt"

type Bipartite struct {
	left  []int
	right []int
	p     [][]int
	q     [][]int
}

func (b *Bipartite) kuhnstep(v int, mem map[int]bool, mt *Matching) bool {
	if mem[v] {
		return false
	}

	mem[v] = true

	for _, u := range b.p[v] {
		if mt.ji[u] == -1 || b.kuhnstep(mt.ji[u], mem, mt) {
			mt.ji[u] = v
			mt.ij[v] = u
			return true
		}
	}

	return false
}

func (b *Bipartite) dfs(v int, mem map[int]bool) [][]int {
	mem[v] = true
	path := [][]int{[]int{v}, []int{}}
	for _, u := range b.p[v] {
		for _, vv := range b.q[u] {
			if seen := mem[vv]; !seen {
				path[1] = append(path[1], u)
				path2 := b.dfs(vv, mem)
				path[0] = append(path[0], path2[0]...)
				path[1] = append(path[1], path2[1]...)
			}
		}
	}

	return path
}

type Matching struct {
	ij []int
	ji []int
}

func (m *Matching) Size() int {
	res := 0
	for _, el := range m.ij {
		if el != -1 {
			res += 1
		}
	}
	return res
}

func (m *Matching) Print() {
	for i, j := range m.ij {
		if j != -1 {
			fmt.Printf("%d -> %d\n", i, j)
		}
	}
}

func main() {
	// c := [][]int{
	// 	{11, 6, 12},
	// 	{12, 4, 6},
	// 	{8, 12, 11},
	// }

	c := [][]int{
		{5, 8, 6, 4},
		{1, 2, 1, 2},
		{3, 2, 8, 5},
		{5, 4, 6, 3},
	}

	n := len(c)

	mt := Matching{ij: make([]int, n), ji: make([]int, n)}
	for i := range mt.ij {
		mt.ij[i] = -1
		mt.ji[i] = -1
	}

	for i, row := range c {
		delta := -1
		for _, el := range row {
			if delta == -1 || delta > el {
				delta = el
			}
		}

		for j := range row {
			c[i][j] -= delta
		}
	}

	for mt.Size() < n {
		// build Bipartite
		b := Bipartite{left: make([]int, n), right: make([]int, n), p: make([][]int, n), q: make([][]int, n)}
		for i := range b.left {
			b.left[i] = -1
		}
		for i := range b.right {
			b.right[i] = -1
		}

		for i, row := range c {
			for j, el := range row {
				if el == 0 {
					b.left[i] = i
					b.p[i] = append(b.p[i], j)
					b.right[j] = j
					b.q[j] = append(b.q[j], i)
				}
			}
		}

		for i, el := range b.left {
			if el == -1 {
				b.left = append(b.left[:i], b.left[i+1:]...)
				b.p = append(b.p[:i], b.p[i+1:]...)
			}
		}
		for j, el := range b.right {
			if el == -1 {
				b.right = append(b.right[:j], b.right[j+1:]...)
				b.q = append(b.q[:j], b.q[j+1:]...)
			}
		}

		mt = Matching{ij: make([]int, n), ji: make([]int, n)}
		for i := range mt.ij {
			mt.ij[i] = -1
			mt.ji[i] = -1
		}
		for i := range b.left {
			mem := make(map[int]bool)
			b.kuhnstep(i, mem, &mt)
		}

		if mt.Size() == n {
			continue
		}

		// run dfs from first free vertex of a left part
		// check what we mark
		// update matrix based on marted left and right marked vertex ( include rows of a left part, exclude columns of a right part)

		freev := -1
		for i, j := range mt.ij {
			if j == -1 {
				freev = i
				break
			}
		}

		mem := make(map[int]bool)
		path := b.dfs(freev, mem)
		for _, i := range path[0] {
			// find min for this row

			delta := -1
			exludeidxs := make(map[int]bool)
			for _, rplus := range path[1] {
				exludeidxs[rplus] = true
			}
			for j, el := range c[i] {
				if exludeidxs[j] {
					continue
				}
				if delta == -1 || delta > el {
					delta = el
				}
			}

			for j := range c[i] {
				if exludeidxs[j] {
					continue
				}
				c[i][j] -= delta
			}
		}
	}

	mt.Print()

}
