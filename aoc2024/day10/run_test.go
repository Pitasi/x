package main

import "testing"

func TestRoots(t *testing.T) {
	input := []byte(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`)

	grid := Parse(input)
	roots := Roots(grid)

	if len(roots) != 9 {
		t.Errorf("Expected 9 roots, got %d", len(roots))
	}
}

func TestScore(t *testing.T) {
	input := []byte(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`)

	grid := Parse(input)

	nines := HikingTrails(grid[0][2], 9)
	score := CountUniques(nines)

	if score != 5 {
		t.Errorf("Expected score 5, got %d", score)
	}
}
