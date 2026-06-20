package main

import (
	"fmt"
	"math/rand"
)

// Edge представляет направленное ребро с пропускной способностью
type Edge struct {
	To  int
	Cap int
}

// Graph — простое представление графа списками смежности
type Graph struct {
	N     int // количество вершин
	Edges [][]Edge
}

// NewGraph создаёт пустой граф на n вершин
func NewGraph(n int) *Graph {
	return &Graph{
		N:     n,
		Edges: make([][]Edge, n),
	}
}

// AddEdge добавляет направленное ребро u -> v с заданной capacity
func (g *Graph) AddEdge(u, v, cap int) {
	g.Edges[u] = append(g.Edges[u], Edge{To: v, Cap: cap})
}

// GenerateRandomGraph создаёт случайный СИЛЬНО связный ориентированный граф
// (из любой вершины достижима любая другая по направленным рёбрам)
// n - количество вершин
// extraPerVertex - сколько дополнительных исходящих рёбер добавить КАЖДОЙ вершине
// maxCap - максимальная пропускная способность ребра (случайная от 1 до maxCap)
func GenerateRandomGraph(n, extraPerVertex, maxCap int, seed int64) *Graph {
	r := rand.New(rand.NewSource(seed))
	g := NewGraph(n)

	// 1. Гарантируем сильную связность: строим случайный гамильтонов цикл
	// perm[0] -> perm[1] -> perm[2] -> ... -> perm[n-1] -> perm[0]
	// Цикл гарантирует, что из любой вершины достижима любая другая
	// (просто идём по кругу), независимо от направления остальных рёбер.
	perm := r.Perm(n)
	for i := 0; i < n; i++ {
		from := perm[i]
		to := perm[(i+1)%n]
		cap := r.Intn(maxCap) + 1
		g.AddEdge(from, to, cap)
	}

	// 2. Каждой вершине добавляем extraPerVertex случайных исходящих рёбер.
	// Это даёт реальное ветвление графа (несколько исходящих рёбер на вершину),
	// а не зависит от случайного попадания при равномерной выборке по всем парам.
	for u := 0; u < n; u++ {
		for k := 0; k < extraPerVertex; k++ {
			v := r.Intn(n)
			if v == u {
				continue
			}
			cap := r.Intn(maxCap) + 1
			g.AddEdge(u, v, cap)
		}
	}

	return g
}

// PickSourceSink выбирает две случайные различные вершины: source и sink
func PickSourceSink(n int, seed int64) (s, t int) {
	r := rand.New(rand.NewSource(seed))
	s = r.Intn(n)
	t = r.Intn(n)
	for t == s {
		t = r.Intn(n)
	}
	return s, t
}

// PrintGraph выводит граф в читаемом виде
func (g *Graph) PrintGraph() {
	fmt.Printf("Граф: %d вершин\n", g.N)
	for u := 0; u < g.N; u++ {
		for _, e := range g.Edges[u] {
			fmt.Printf("  %d -> %d (cap=%d)\n", u, e.To, e.Cap)
		}
	}
}
