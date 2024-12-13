package main

import (
	"testing"
)

func TestStar1SingleGame(t *testing.T) {
	game := Game{
		A:     Position{94, 34},
		B:     Position{22, 67},
		Prize: Position{8400, 5400},
	}

	actual := Star1(game)
	expected := 280
	if actual != expected {
		t.Fatalf("expected %d, got %d", expected, actual)
	}
}

func TestStar1SingleGameImpossible(t *testing.T) {
	game := Game{
		A:     Position{26, 66},
		B:     Position{67, 21},
		Prize: Position{12748, 12176},
	}

	actual := Star1(game)
	expected := 0
	if actual != expected {
		t.Fatalf("expected %d, got %d", expected, actual)
	}
}

func TestStar1SingleGameHigh(t *testing.T) {
	game := Game{
		A:     Position{26, 66},
		B:     Position{67, 21},
		Prize: Position{10000000012748, 10000000012176},
	}

	actual := Star1(game)
	expected := 0
	if actual != expected {
		t.Fatalf("expected %d, got %d", expected, actual)
	}
}

func TestStar1(t *testing.T) {
	input := []byte(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`)

	games, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 480
	var actual int
	for _, g := range games {
		actual += Star1(g)
	}
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}
