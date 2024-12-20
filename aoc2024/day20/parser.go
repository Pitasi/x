package main

import (
	"fmt"
	"strings"
)

type Cell struct {
	v             byte
	remainingTime int
}

func (c *Cell) SetRemainingTime(t int) {
	if c.v != '.' {
		panic(fmt.Sprintf("Can't set remaining time on non-empty cell. Value: %c", c.v))
	}
	c.remainingTime = t
}

func (c *Cell) IsWall() bool {
	return c == nil || c.v == '#'
}

type Grid struct {
	m          [][]*Cell
	track      []Vec
	start, end Vec
}

func (g *Grid) String() string {
	var sb strings.Builder
	for y, row := range g.m {
		for x, cell := range row {
			if x == g.start.X && y == g.start.Y {
				sb.WriteString("S")
			} else if x == g.end.X && y == g.end.Y {
				sb.WriteString("E")
			} else {
				sb.WriteByte(cell.v)
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (g *Grid) PrintVec(v1, v2 Vec) {
	var sb strings.Builder
	for y, row := range g.m {
		for x, cell := range row {
			if x == v1.X && y == v1.Y {
				sb.WriteString("1")
			} else if x == v2.X && y == v2.Y {
				sb.WriteString("2")
			} else if x == g.start.X && y == g.start.Y {
				sb.WriteString("S")
			} else if x == g.end.X && y == g.end.Y {
				sb.WriteString("E")
			} else {
				sb.WriteByte(cell.v)
			}
		}
		sb.WriteByte('\n')
	}
	fmt.Print(sb.String())
}

func (g *Grid) Width() int {
	return len(g.m[0])
}

func (g *Grid) Height() int {
	return len(g.m)
}

func (g *Grid) Cell(x, y int) *Cell {
	if x < 0 || x >= g.Width() {
		return nil
	}
	if y < 0 || y >= g.Height() {
		return nil
	}
	return g.m[y][x]
}

func (g *Grid) computeTrack() {
	var (
		cur     = g.start
		curTime = 1
	)

	for {
		g.track = append(g.track, cur)
		c := g.Cell(cur.X, cur.Y)
		c.SetRemainingTime(curTime)
		curTime++

		up := g.Cell(cur.X, cur.Y-1)
		down := g.Cell(cur.X, cur.Y+1)
		left := g.Cell(cur.X-1, cur.Y)
		right := g.Cell(cur.X+1, cur.Y)

		if !up.IsWall() && up.remainingTime == 0 {
			cur = Vec{cur.X, cur.Y - 1}
			continue
		} else if !down.IsWall() && down.remainingTime == 0 {
			cur = Vec{cur.X, cur.Y + 1}
			continue
		} else if !left.IsWall() && left.remainingTime == 0 {
			cur = Vec{cur.X - 1, cur.Y}
			continue
		} else if !right.IsWall() && right.remainingTime == 0 {
			cur = Vec{cur.X + 1, cur.Y}
			continue
		} else {
			break
		}
	}
}

type Vec struct {
	X, Y int
}

func Parse(input []byte) *Grid {
	var rows [][]*Cell

	var (
		pos    int
		curRow []*Cell
		S, E   Vec
	)

	for _, b := range input {
		switch b {
		case '#':
			curRow = append(curRow, &Cell{b, 0})
		case '.':
			curRow = append(curRow, &Cell{b, 0})
		case 'S':
			S = Vec{len(curRow), len(rows)}
			curRow = append(curRow, &Cell{'.', 0})
		case 'E':
			E = Vec{len(curRow), len(rows)}
			curRow = append(curRow, &Cell{'.', 0})
		case '\n':
			rows = append(rows, curRow)
			curRow = nil
		default:
			panic("Unexpected byte: " + string(b))
		}
		pos++
	}
	if len(curRow) > 0 {
		rows = append(rows, curRow)
	}

	g := &Grid{
		m:     rows,
		start: S,
		end:   E,
	}
	g.computeTrack()

	return g
}
