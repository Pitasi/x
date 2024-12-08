package main

import "testing"

func TestRun(t *testing.T) {
	m := Map{
		grid: [][]Cell{
			{e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), a('0'), e(), e(), e()},
			{e(), e(), e(), e(), e(), a('0'), e(), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), a('0'), e(), e(), e(), e()},
			{e(), e(), e(), e(), a('0'), e(), e(), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), a('A'), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), a('A'), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), e(), a('A'), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e()},
			{e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e(), e()},
		},
	}

	count := Run(m)
	expected := 34
	if count != expected {
		t.Errorf("expected %d, got %d", expected, count)
		printMap(m)
	}
}

func e() Cell {
	return Cell{}
}

func a(c byte) Cell {
	return Cell{antenna: c}
}
