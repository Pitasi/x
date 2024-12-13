package main

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`Button A: X+11, Y+73
Button B: X+65, Y+17
Prize: X=18133, Y=4639

Button A: X+49, Y+13
Button B: X+24, Y+79
Prize: X=6664, Y=948`)

	games, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(games) != 2 {
		t.Errorf("expected 2 games, got %d", len(games))
	}

	t.Logf("%+v", games[0])
	t.Logf("%+v", games[1])
}
