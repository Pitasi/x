package main

import "testing"

func TestParse(t *testing.T) {
	input := `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`

	g := Parse([]byte(input))

	width := 15
	height := 15

	if g.Width() != width {
		t.Errorf("Expected width %d, got %d", width, g.Width())
	}

	if g.Height() != height {
		t.Errorf("Expected height %d, got %d", height, g.Height())
	}

	if g.Cell(0, 0).v != '#' {
		t.Errorf("Expected cell at 0,0 to be #, got %c", g.Cell(0, 0))
	}
}
