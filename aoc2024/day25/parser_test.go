package main

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`#####
.####
.####
.####
.#.#.
.#...
.....`)

	expected := Key{
		a: 0, b: 5, c: 3, d: 4, e: 3,
		isLock: true,
	}

	actual := Parse(input)
	if len(actual) != 1 {
		t.Errorf("expected %d keys, got %d", 1, len(actual))
	}
	if actual[0] != expected {
		t.Errorf("expected %v, got %v", expected, actual[0])
	}
}
