package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	g := parseGrid(input)

	// star1(g)
	star2(g)
}

func star1(g grid) {
	var count int
	for h := 0; h < g.height; h++ {
		for x := 0; x < g.width; x++ {
			c := g.get(x, h)
			if c == 'S' {
				g.set(x, h, '|')
			}
			if c == '^' && g.get(x, h-1) == '|' {
				count++
				g.set(x-1, h, '|')
				g.set(x+1, h, '|')
			} else if g.get(x, h-1) == '|' {
				g.set(x, h, '|')
			}
		}
	}
	fmt.Println(count)
}

func star2(g grid) {
	nodeG := newNodeGrid(g.width, g.height)
	s := start(g)
	s.findChildren(nodeG, g, s.x)
	fmt.Println(s.timelines())
}

func start(g grid) *node {
	for x := 0; x < g.width; x++ {
		if g.get(x, 0) == 'S' {
			return &node{x: x, y: 0}
		}
	}
	panic("start not found")
}

type grid struct {
	rows   [][]byte
	height int
	width  int
}

func (g grid) get(x, y int) byte {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return 0
	}
	return g.rows[y][x]
}

func (g grid) set(x, y int, c byte) {
	g.rows[y][x] = c
}

func (g grid) String() string {
	return string(bytes.Join(g.rows, []byte("\n")))
}

func parseGrid(in []byte) grid {
	lines := bytes.Split(in, []byte("\n"))

	rows := make([][]byte, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		rows = append(rows, line)
	}

	height := len(rows)
	var width int
	if height > 0 {
		width = len(rows[0])
	}

	// input check: all rows are same width
	for _, r := range rows {
		if len(r) != width {
			panic(fmt.Errorf("mismtched width: expected %d got %d", width, len(r)))
		}
	}

	return grid{rows, height, width}
}

type node struct {
	x        int
	y        int
	children []*node
	tl       *int
}

func (n *node) timelines() int {
	if n.tl != nil {
		return *n.tl
	}

	var r int
	if len(n.children) == 0 {
		r = 2
	}
	if len(n.children) == 1 {
		r = n.children[0].timelines()
		if n.y > 0 {
			r++
		}
	}
	if len(n.children) == 2 {
		r = n.children[0].timelines() + n.children[1].timelines()
	}

	n.tl = &r
	return r
}

func (n *node) findChildren(ng *nodeGrid, g grid, x int) {
	for y := n.y + 1; y < g.height; y++ {
		sub := g.get(x, y)
		if sub == '.' {
			continue
		}
		if sub == '^' {
			child := ng.get(x, y)
			if child == nil {
				child = &node{
					x: x,
					y: y,
				}
				child.findChildren(ng, g, x-1)
				child.findChildren(ng, g, x+1)
				ng.set(x, y, child)
			}
			n.children = append(n.children, child)
			break
		}
	}
}

type nodeGrid struct {
	rows   [][]*node
	height int
	width  int
}

func newNodeGrid(w, h int) *nodeGrid {
	rows := make([][]*node, h)
	for i := range h {
		rows[i] = make([]*node, w)
	}
	return &nodeGrid{rows: rows, height: h, width: w}
}

func (g nodeGrid) get(x, y int) *node {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return nil
	}
	return g.rows[y][x]
}

func (g nodeGrid) set(x, y int, n *node) {
	g.rows[y][x] = n
}
