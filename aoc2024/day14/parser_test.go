package main

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`)
	robots := Parse(input)
	expected := []Robot{
		{P: v(0, 4), V: v(3, -3)},
		{P: v(6, 3), V: v(-1, -3)},
		{P: v(10, 3), V: v(-1, 2)},
		{P: v(2, 0), V: v(2, -1)},
		{P: v(0, 0), V: v(1, 3)},
		{P: v(3, 0), V: v(-2, -2)},
		{P: v(7, 6), V: v(-1, -3)},
		{P: v(3, 0), V: v(-1, -2)},
		{P: v(9, 3), V: v(2, 3)},
		{P: v(7, 3), V: v(-1, 2)},
		{P: v(2, 4), V: v(2, -3)},
		{P: v(9, 5), V: v(-3, -3)},
	}
	if len(robots) != len(expected) {
		t.Errorf("expected %d robots, got %d", len(expected), len(robots))
	}
	for i := range robots {
		if robots[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], robots[i])
		}
	}
}
