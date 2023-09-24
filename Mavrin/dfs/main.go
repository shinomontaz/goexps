package main

import "fmt"

type Node struct {
	id  int
	out []*Node
}

type Graph struct {
	nodes []*Node
}

func (g *Graph) Add(node *Node) {
	g.nodes = append(g.nodes, node)
}

// u - start node
func (g *Graph) Dfs(u *Node, mem map[*Node]bool) []int {
	mem[u] = true
	path := []int{u.id}
	for _, v := range u.out {
		if seen := mem[v]; !seen {
			path = append(path, g.Dfs(v, mem)...)
		}
	}

	return path
}

func main() {
	g := Graph{
		nodes: make([]*Node, 0),
	}
	n1 := &Node{id: 1}
	n2 := &Node{id: 2}
	n3 := &Node{id: 3}
	n4 := &Node{id: 4}
	n5 := &Node{id: 5}

	n1.out = []*Node{n2, n3}
	n2.out = []*Node{n4, n5}
	n3.out = []*Node{n5}

	g.Add(n1)
	g.Add(n2)
	g.Add(n3)
	g.Add(n4)
	g.Add(n5)

	mem := make(map[*Node]bool)
	path := g.Dfs(n1, mem)

	fmt.Println(path)
}
