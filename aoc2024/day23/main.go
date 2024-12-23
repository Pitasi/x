package main

import (
	_ "embed"
	"fmt"
	"iter"
	"slices"
	"strings"
)

//go:embed input.txt
var input []byte

func main() {
	g := Parse(input)
	fmt.Println(star2(g))
}

func star2(g Graph) string {
	var maxVisited map[string]*Node
	for _, start := range g {
		visited := make(map[string]*Node)
		visited[start.name] = start
		for _, end := range g {
			if start == end {
				continue
			}

			connected := true
			for _, v := range visited {
				if !end.IsConnectedTo(v) {
					connected = false
					break
				}
			}

			if connected {
				visited[end.name] = end
			}
		}

		if len(visited) > len(maxVisited) {
			maxVisited = visited
		}
	}
	names := make([]string, 0, len(maxVisited))
	for n := range maxVisited {
		names = append(names, n)
	}
	return identifierNames(names...)
}

func star1(g Graph) int {
	s := make(map[string]struct{})
	for _, n := range g {
		if n.name[0] != 't' {
			continue
		}
		for n1, n2 := range n.ConnectionsTuples() {
			if n1.IsConnectedTo(n2) {
				s[identifier(n, n1, n2)] = struct{}{}
			}
		}
	}
	return len(s)
}

func identifier(nodes ...*Node) string {
	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.name
	}
	return identifierNames(names...)
}

func identifierNames(names ...string) string {
	slices.Sort(names)
	return strings.Join(names, ",")
}

type Graph map[string]*Node

func newGraph() Graph {
	return make(Graph)
}

func (g Graph) AddNode(name string) *Node {
	if n, ok := g[name]; ok {
		return n
	}
	node := &Node{
		name:        name,
		connections: make([]*Node, 0),
	}
	g[name] = node
	return node
}

func (g Graph) AddEdge(from, to string) {
	g[from].AddConnection(g[to])
	g[to].AddConnection(g[from])
}

type Node struct {
	name        string
	connections []*Node
}

func (n *Node) AddConnection(other *Node) {
	n.connections = append(n.connections, other)
}

func (n *Node) ConnectionsTuples() iter.Seq2[*Node, *Node] {
	return func(yield func(*Node, *Node) bool) {
		for i := 0; i < len(n.connections)-1; i++ {
			for j := i + 1; j < len(n.connections); j++ {
				if !yield(n.connections[i], n.connections[j]) {
					return
				}
			}
		}
	}
}

func (n *Node) IsConnectedTo(other *Node) bool {
	for _, c := range n.connections {
		if c.name == other.name {
			return true
		}
	}
	return false
}
