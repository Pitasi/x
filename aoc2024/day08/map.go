package main

type Map struct {
	grid [][]Cell
}

func (m Map) GetCell(pos Position) (*Cell, bool) {
	if pos.y < 0 || pos.y >= len(m.grid) {
		return nil, false
	}
	if pos.x < 0 || pos.x >= len(m.grid[pos.y]) {
		return nil, false
	}

	return &m.grid[pos.y][pos.x], true
}

type Position struct {
	x int
	y int
}

func (p Position) Up() Position {
	return Position{x: p.x, y: p.y - 1}
}

func (p Position) Down() Position {
	return Position{x: p.x, y: p.y + 1}
}

func (p Position) Left() Position {
	return Position{x: p.x - 1, y: p.y}
}

func (p Position) Right() Position {
	return Position{x: p.x + 1, y: p.y}
}

func (p Position) Antinode(p2 Position, boundX, boundY int) []Position {
	dx := p2.x - p.x
	dy := p2.y - p.y

	var (
		nodes []Position
		x     = p2.x
		y     = p2.y
	)
	for x >= 0 && y >= 0 && x < boundX && y < boundY {
		nodes = append(nodes, Position{x: x, y: y})
		x += dx
		y += dy
	}

	return nodes
}

type Cell struct {
	antenna  byte
	antinode bool
}

func (c *Cell) MarkAntinode() { c.antinode = true }
