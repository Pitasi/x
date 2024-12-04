package main

import (
	"bytes"
	_ "embed"
	"fmt"
)

type Direction int

const (
	DOWN Direction = iota
	RIGHT
	LEFT
	UP
	UP_RIGHT
	UP_LEFT
	DOWN_RIGHT
	DOWN_LEFT
)

type Grid struct {
	rows [][]byte
}

func (g *Grid) Get4(x, y int, dir Direction) []byte {
	var (
		xs []int
		ys []int
	)

	switch dir {
	case UP:
		xs = []int{x, x, x, x}
		ys = []int{y, y - 1, y - 2, y - 3}
	case DOWN:
		xs = []int{x, x, x, x}
		ys = []int{y, y + 1, y + 2, y + 3}
	case RIGHT:
		xs = []int{x, x + 1, x + 2, x + 3}
		ys = []int{y, y, y, y}
	case LEFT:
		xs = []int{x, x - 1, x - 2, x - 3}
		ys = []int{y, y, y, y}
	case DOWN_RIGHT:
		xs = []int{x, x + 1, x + 2, x + 3}
		ys = []int{y, y + 1, y + 2, y + 3}
	case DOWN_LEFT:
		xs = []int{x, x - 1, x - 2, x - 3}
		ys = []int{y, y + 1, y + 2, y + 3}
	case UP_RIGHT:
		xs = []int{x, x + 1, x + 2, x + 3}
		ys = []int{y, y - 1, y - 2, y - 3}
	case UP_LEFT:
		xs = []int{x, x - 1, x - 2, x - 3}
		ys = []int{y, y - 1, y - 2, y - 3}
	default:
		panic("Grid.Get4: invalid direction")
	}

	res := make([]byte, 4)
	for i := 0; i < 4; i++ {
		v, ok := g.get(xs[i], ys[i])
		if !ok {
			return nil
		}
		res[i] = v
	}

	return res
}

func (g *Grid) GetX(x, y int) [][]byte {
	xs := []int{x - 1, x, x + 1}
	ys := []int{y - 1, y, y + 1}

	var res [][]byte
	for y := range len(ys) {
		var row []byte
		for x := range len(xs) {
			v, ok := g.get(xs[x], ys[y])
			if !ok {
				return nil
			}
			row = append(row, v)
		}
		res = append(res, row)
	}

	return res
}

func (g *Grid) get(x, y int) (byte, bool) {
	if y < 0 || y > len(g.rows)-1 {
		return 0, false
	}
	row := g.rows[y]
	if x < 0 || x > len(row)-1 {
		return 0, false
	}
	return row[x], true
}

func (g *Grid) Rows() int { return len(g.rows) }

func (g *Grid) Cols() int { return len(g.rows[0]) }

func parseGrid(bz []byte) Grid {
	var currentRow []byte
	rows := make([][]byte, 0)
	for _, b := range bz {
		if b == '\n' {
			rows = append(rows, currentRow)
			currentRow = make([]byte, 0)
		} else {
			currentRow = append(currentRow, b)
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}
	return Grid{
		rows: rows,
	}
}

//go:embed input.txt
var input []byte

func main() {
	g := parseGrid(input)
	count := CountXmas(g)
	fmt.Println(count)
	fmt.Println(CountXs(g))
}

var directions = []Direction{
	UP,
	DOWN,
	RIGHT,
	LEFT,
	UP_RIGHT,
	UP_LEFT,
	DOWN_RIGHT,
	DOWN_LEFT,
}

func CountXmas(g Grid) int {
	var count int
	for y := range g.Rows() {
		for x := range g.Cols() {
			for _, dir := range directions {
				b := g.Get4(x, y, dir)
				if isXmas(b) {
					count++
				}
			}
		}
	}
	return count
}

func CountXs(g Grid) int {
	count := 0
	for y := range g.Rows() {
		for x := range g.Cols() {
			b := g.GetX(x, y)
			if isValidX(b) {
				count++
			}
		}
	}
	return count
}

func isXmas(b []byte) bool {
	return bytes.Equal(b, []byte("XMAS"))
}

func isValidX(b [][]byte) bool {
	if len(b) != 3 {
		return false
	}

	center := b[1][1]
	topLeft := b[0][0]
	topRight := b[0][2]
	bottomLeft := b[2][0]
	bottomRight := b[2][2]

	if center != 'A' {
		return false
	}

	return (topLeft == 'M' && topRight == 'M' && bottomLeft == 'S' && bottomRight == 'S' ||
		topLeft == 'S' && topRight == 'S' && bottomLeft == 'M' && bottomRight == 'M' ||
		topLeft == 'M' && topRight == 'S' && bottomLeft == 'M' && bottomRight == 'S' ||
		topLeft == 'S' && topRight == 'M' && bottomLeft == 'S' && bottomRight == 'M')
}
