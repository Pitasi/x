package main

import (
	"fmt"
	"reflect"
	"testing"
)

var testInput = `###############
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

func TestCheatCandidates(t *testing.T) {
	g := Parse([]byte(testInput))
	hCand, vCand := CheatV1Candidates(g)
	expected := 44
	if len(hCand)+len(vCand) != expected {
		t.Errorf("Expected %d candidates, got %d", expected, len(hCand)+len(vCand))
	}
}

func TestCheatSavingH(t *testing.T) {
	g := Parse([]byte(testInput))

	wall := Vec{8, 1}
	expected := 12
	got := CheatV1SavingH(g, wall)
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestExample(t *testing.T) {
	g := Parse([]byte(testInput))
	candsH, candsV := CheatV1Candidates(g)

	savings := make(map[int]int)
	for _, c := range candsH {
		s := CheatV1SavingH(g, c)
		savings[s]++
	}
	for _, c := range candsV {
		s := CheatV1SavingV(g, c)
		savings[s]++
	}

	expected := map[int]int{
		2:  14,
		4:  14,
		6:  2,
		8:  4,
		10: 2,
		12: 3,
		20: 1,
		36: 1,
		38: 1,
		40: 1,
		64: 1,
	}
	if !reflect.DeepEqual(savings, expected) {
		t.Errorf("Expected %v, got %v", expected, savings)
	}
}

func TestExampleV2(t *testing.T) {
	g := Parse([]byte(testInput))
	cands := CheatV2Candidates(g, 0, 2)

	savings := make(map[int]int)
	for _, c := range cands {
		s := CheatV2Saving(g, c)
		if s > 0 {
			savings[s]++
		}
	}

	expected := map[int]int{
		2:  14,
		4:  14,
		6:  2,
		8:  4,
		10: 2,
		12: 3,
		20: 1,
		36: 1,
		38: 1,
		40: 1,
		64: 1,
	}
	if !reflect.DeepEqual(savings, expected) {
		t.Errorf("Expected %v, got %v", expected, savings)
	}
}

func TestCheatV2Saving(t *testing.T) {
	g := Parse([]byte(testInput))
	c := CheatV2{Vec{1, 3}, Vec{3, 7}, 6}
	expected := 76
	got := CheatV2Saving(g, c)
	if got != expected {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func TestExample2(t *testing.T) {
	g := Parse([]byte(testInput))
	cands := CheatV2Candidates(g, 50, 20)

	savings := make(map[int]int)
	for _, c := range cands {
		s := CheatV2Saving(g, c)
		savings[s]++

		if s == 70 {
			fmt.Println("Candidate:", c, "saving:", s)
			//g.PrintVec(c.Start, c.End)
		}
	}

	expected := map[int]int{
		50: 32,
		52: 31,
		54: 29,
		56: 39,
		58: 25,
		60: 23,
		62: 20,
		64: 19,
		66: 12,
		68: 14,
		70: 12,
		72: 22,
		74: 4,
		76: 3,
	}
	if !reflect.DeepEqual(savings, expected) {
		t.Errorf("Expected %v, got %v", expected, savings)
	}
}

func TestGetCheatCost(t *testing.T) {
	g := Parse([]byte(testInput))
	start := Vec{1, 3}
	end := Vec{3, 7}
	maxCheatTime := 20

	cost, found := getCheatCost(g, start, end, maxCheatTime)
	if !found {
		t.Fatalf("Expected to find a path from %v to %v", start, end)
	}
	expected := 6
	if cost != expected {
		t.Errorf("Expected cost %d, got %d", expected, cost)
	}
}

func TestGetCheatCost2(t *testing.T) {
	g := Parse([]byte(testInput))
	start := Vec{1, 3}
	end := Vec{4, 7}
	maxCheatTime := 20

	cost, found := getCheatCost(g, start, end, maxCheatTime)
	if !found {
		t.Fatalf("Expected to find a path from %v to %v", start, end)
	}
	expected := 13
	if cost != expected {
		t.Errorf("Expected cost %d, got %d", expected, cost)
	}
}

func TestSmall(t *testing.T) {
	g := Parse([]byte(`######
#....#
#.##.#
#.##.#
#.##.#
#.##.#
#S##E#
######`))
	cands := CheatV2Candidates(g, 0, 2)

	t.Log(cands)
}
