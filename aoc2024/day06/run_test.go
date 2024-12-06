package main

import "testing"

func TestRun(t *testing.T) {
	m := Map{
		grid: [][]Cell{
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}},
			{&Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}},
		},
		guard: Guard{
			position:  Position{x: 4, y: 6},
			direction: Up,
		},
	}

	count, _ := Run(m)

	if count != 41 {
		t.Errorf("expected 41, got %d", count)
	}
}

func TestCountLoops(t *testing.T) {
	m := Map{
		grid: [][]Cell{
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}},
			{&Wall{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}},
			{&Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Empty{}, &Wall{}, &Empty{}, &Empty{}, &Empty{}},
		},
		guard: Guard{
			position:  Position{x: 4, y: 6},
			direction: Up,
		},
	}

	Run(m)
	count := CountLoops(m)

	if count != 6 {
		t.Errorf("expected 6, got %d", count)
	}
}
