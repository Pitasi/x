package main

import "testing"

func TestParser(t *testing.T) {
	input := []byte("2333133121414131402")
	expected := &Disk{
		blocks: []*Block{
			{data: 0},
			{data: 0},
			nil,
			nil,
			nil,
			{data: 1},
			{data: 1},
			{data: 1},
			nil,
			nil,
			nil,
			{data: 2},
			nil,
			nil,
			nil,
			{data: 3},
			{data: 3},
			{data: 3},
			nil,
			{data: 4},
			{data: 4},
			nil,
			{data: 5},
			{data: 5},
			{data: 5},
			{data: 5},
			nil,
			{data: 6},
			{data: 6},
			{data: 6},
			{data: 6},
			nil,
			{data: 7},
			{data: 7},
			{data: 7},
			nil,
			{data: 8},
			{data: 8},
			{data: 8},
			{data: 8},
			{data: 9},
			{data: 9},
		},
	}

	actual := Parse(input)
	if actual.checksum() != expected.checksum() {
		t.Errorf("expected checksum %d, got %d", expected.checksum(), actual.checksum())
	}
}
