package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	g := &graph{nodes: make(map[string]*node)}

	lines := strings.SplitSeq(string(input), "\n")
	for l := range lines {
		if len(l) == 0 {
			continue
		}

		name := l[:3]
		node := g.FindOrCreate(name)

		conns := strings.SplitSeq(l[5:], " ")
		for c := range conns {
			target := g.FindOrCreate(c)
			node.Connect(target)
		}
	}

	// star1(g)
	star2(g)
}

func star1(g *graph) {
	you := g.Find("you")
	out := g.Find("out")
	count := visit(you, out)
	fmt.Println(count)
}

func star2(g *graph) {
	svr := g.Find("svr")
	out := g.Find("out")
	dac := g.Find("dac")
	fft := g.Find("fft")

	// c1 := visit(svr, dac)
	// c2 := visit(dac, fft) // this returned 0 so all this group can be excluded
	// c3 := visit(fft, out)

	a1 := visit2(svr, fft)
	g.ResetPaths()
	a2 := visit2(fft, dac)
	g.ResetPaths()
	a3 := visit2(dac, out)
	g.ResetPaths()

	fmt.Println(a1 * a2 * a3)
}

func visit(start, target *node) int {
	var count int

	for c := range start.connections {
		if c == target {
			count++
			continue
		}
		count += visit(c, target)
	}

	return count
}

func visit2(start, target *node) int {
	if start.paths != nil {
		return *start.paths
	}

	var c int
	start.paths = &c

	for c := range start.connections {
		if c == target {
			(*start.paths)++
			continue
		}
		(*start.paths) += visit2(c, target)
	}

	return (*start.paths)
}

type graph struct {
	nodes map[string]*node
}

func (g *graph) Find(name string) *node {
	return g.nodes[name]
}

func (g *graph) ResetPaths() {
	for _, n := range g.nodes {
		n.paths = nil
	}
}

func (g *graph) FindOrCreate(name string) *node {
	if n, found := g.nodes[name]; found {
		return n
	}

	n := &node{
		name:        name,
		connections: make(map[*node]struct{}),
	}

	g.nodes[name] = n

	return n
}

type node struct {
	name        string
	connections map[*node]struct{}
	paths       *int
}

func (n *node) Connect(other *node) {
	n.connections[other] = struct{}{}
}
