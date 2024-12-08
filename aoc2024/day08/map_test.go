package main

import "testing"

func TestAntinode(t *testing.T) {
	p1 := Position{x: 4, y: 3}
	p2 := Position{x: 5, y: 5}

	a1 := p1.Antinode(p2, 10, 10)
	expected := Position{x: 6, y: 7}
	if a1[1] != expected {
		t.Errorf("got %v, expected %v", a1, expected)
	}

	a2 := p2.Antinode(p1, 10, 10)
	expected = Position{x: 3, y: 1}
	if a2[1] != expected {
		t.Errorf("got %v, expected %v", a2, expected)
	}
}
