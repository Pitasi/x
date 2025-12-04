package main

import (
	"bytes"
	"fmt"
	"iter"
	"os"
)

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	g := parseGrid(input)

	star1(g)
	star2(g)
}

func star1(g grid) {
	var count int
	for x, y := range g.iterCoords() {
		if g.get(x, y) != '@' {
			continue
		}

		adj := g.getAdjiacent(x, y)
		if lessThan(adj, '@', 4) {
			count++
		}
	}
	fmt.Println(count)
}

func star2(g grid) {
	var count int
	for {
		// mark
		var markX []int
		var markY []int
		for x, y := range g.iterCoords() {
			if g.get(x, y) != '@' {
				continue
			}

			adj := g.getAdjiacent(x, y)
			if lessThan(adj, '@', 4) {
				markX = append(markX, x)
				markY = append(markY, y)
			}
		}

		// sweep
		if len(markX) == 0 {
			break
		}

		count += len(markX)
		for i := 0; i < len(markX); i++ {
			g.set(markX[i], markY[i], '.')
		}
	}
	fmt.Println(count)
}

func lessThan(adj [8]byte, c byte, n int) bool {
	var count int
	for _, a := range adj {
		if a == c {
			count++
		}
	}
	return count < n
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

func (g grid) getAdjiacent(x, y int) [8]byte {
	return [8]byte{
		g.get(x-1, y-1),
		g.get(x, y-1),
		g.get(x+1, y-1),
		g.get(x-1, y),
		g.get(x+1, y),
		g.get(x-1, y+1),
		g.get(x, y+1),
		g.get(x+1, y+1),
	}
}

func (g grid) iterCoords() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for y := range g.height {
			for x := range g.width {
				if !yield(x, y) {
					return
				}
			}
		}

	}
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
