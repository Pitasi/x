package main

import "testing"

func TestIsRect(t *testing.T) {
	grid := [][]rune{
		{'X', 'X', 'X', 'X', 'X'},
		{'X', 'X', 'X', 'X', 'X'},
		{'X', 'X', 'X', ' ', 'X'},
		{'X', ' ', ' ', ' ', 'X'},
		{'X', 'X', 'X', 'X', 'X'},
	}
	tests := []struct {
		x, y, w, h int
		want       bool
	}{
		{0, 0, 3, 3, true},
		{1, 0, 3, 3, false},
		{0, 1, 3, 3, false},
		{1, 1, 2, 2, true},
	}

	for _, test := range tests {
		got := isRect(grid, test.x, test.y, test.w, test.h)
		if got != test.want {
			t.Errorf("isRect(%d, %d, %d, %d) = %v, want %v", test.x, test.y, test.w, test.h, got, test.want)
		}
	}
}
