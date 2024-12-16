package main

import "testing"

func TestRun1(t *testing.T) {
	input := `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`

	g, startX, startY, goalX, goalY := Parse([]byte(input))

	score, count := Run(g, startX, startY, goalX, goalY)

	if score != 7036 {
		t.Errorf("expected score to be 11048, got %d", score)
	}

	if count != 45 {
		t.Errorf("expected count to be 64, got %d", count)
	}
}

func TestRun(t *testing.T) {
	input := `#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################`

	g, startX, startY, goalX, goalY := Parse([]byte(input))

	score, count := Run(g, startX, startY, goalX, goalY)

	if score != 11048 {
		t.Errorf("expected score to be 11048, got %d", score)
	}

	if count != 64 {
		t.Errorf("expected count to be 64, got %d", count)
	}
}

func TestParse(t *testing.T) {
	input := `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`

	g, startX, startY, goalX, goalY := Parse([]byte(input))

	expectedStartX := 1
	expectedStartY := 13

	if startX != expectedStartX {
		t.Errorf("expected startX to be %d, got %d", expectedStartX, startX)
	}

	if startY != expectedStartY {
		t.Errorf("expected startY to be %d, got %d", expectedStartY, startY)
	}

	expectedGoalX := 13
	expectedGoalY := 1

	if goalX != expectedGoalX {
		t.Errorf("expected goalX to be %d, got %d", expectedGoalX, goalX)
	}

	if goalY != expectedGoalY {
		t.Errorf("expected goalY to be %d, got %d", expectedGoalY, goalY)
	}

	expected := `###############
#.......#.....#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#...#.....#...#
###############
`

	if g.String() != expected {
		t.Errorf("expected g to be %s, got %s", expected, g.String())
	}
}
