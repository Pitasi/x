package main

import (
	"fmt"
	"strings"
)

type Map struct {
	grid  [][]Cell
	guard Guard
}

func (m Map) GetCell(pos Position) Cell {
	if pos.y < 0 || pos.y >= len(m.grid) {
		return nil
	}
	if pos.x < 0 || pos.x >= len(m.grid[pos.y]) {
		return nil
	}

	return m.grid[pos.y][pos.x]
}

func (m Map) With(x, y int, cell Cell) Map {
	rows := make([][]Cell, len(m.grid))
	for i, row := range m.grid {
		rows[i] = make([]Cell, len(row))
		for j, cell := range row {
			rows[i][j] = cell.Clone()
		}
	}
	rows[y][x] = cell

	return Map{
		grid:  rows,
		guard: m.guard,
	}
}

func (m Map) String() string {
	var b strings.Builder

	for y, row := range m.grid {
		for x, cell := range row {
			if x == m.guard.position.x && y == m.guard.position.y {
				b.WriteString("X")
				continue
			}
			b.WriteString(cell.String())
		}
		b.WriteString("\n")
	}

	return b.String()
}

type Guard struct {
	position  Position
	direction Direction
}

func (g *Guard) SetPosition(p Position) {
	g.position = p
}

func (g *Guard) SetDirection(d Direction) {
	g.direction = d
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

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

type Cell interface {
	fmt.Stringer

	Clone() Cell
}

func IsObstacle(c Cell) bool {
	_, ok := c.(Obstacle)
	return ok
}

type Obstacle interface {
	obstacle()
}

type Marker interface {
	Mark(dir Direction)
	IsMarked() bool
	IsMarkedFor(dir Direction) bool
}

type Wall struct{}

func (w Wall) obstacle() {}

func (w Wall) String() string { return "#" }

func (w Wall) Clone() Cell { return &w }

var _ Obstacle = (*Wall)(nil)

type Empty struct {
	marks map[Direction]struct{}
}

var _ (Marker) = (*Empty)(nil)

func (e Empty) String() string { return "." }

func (e *Empty) Mark(dir Direction) {
	if e.marks == nil {
		e.marks = make(map[Direction]struct{})
	}
	e.marks[dir] = struct{}{}
}

func (e *Empty) IsMarked() bool {
	return e.marks != nil
}

func (e *Empty) IsMarkedFor(dir Direction) bool {
	if e.marks == nil {
		return false
	}
	_, ok := e.marks[dir]
	return ok
}

func (e *Empty) Clone() Cell { return &Empty{} }
