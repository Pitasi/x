package main

import "testing"

func TestStep(t *testing.T) {
	r := Robot{
		P: v(2, 4),
		V: v(2, -3),
	}
	mapSize := v(11, 7)

	r2 := Step(5, r, mapSize)

	expected := Robot{
		P: v(1, 3),
		V: v(2, -3),
	}

	if r2 != expected {
		t.Errorf("expected %v, got %v", expected, r2)
	}
}
