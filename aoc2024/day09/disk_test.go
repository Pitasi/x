package main

import (
	"testing"
)

func TestCompact(t *testing.T) {
	disk := &Disk{
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

	CompactByBlock(disk)

	expected := []*Block{
		{data: 0},
		{data: 0},
		{data: 9},
		{data: 9},
		{data: 8},
		{data: 1},
		{data: 1},
		{data: 1},
		{data: 8},
		{data: 8},
		{data: 8},
		{data: 2},
		{data: 7},
		{data: 7},
		{data: 7},
		{data: 3},
		{data: 3},
		{data: 3},
		{data: 6},
		{data: 4},
		{data: 4},
		{data: 6},
		{data: 5},
		{data: 5},
		{data: 5},
		{data: 5},
		{data: 6},
		{data: 6},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}

	for i := range expected {
		gotBlock := disk.blocks[i]
		expectedBlock := expected[i]

		// both nil
		if gotBlock == nil && expectedBlock == nil {
			continue
		}
		// got nil, expected not nil
		if gotBlock == nil && expectedBlock != nil {
			t.Fatalf("block[%d]: expected %d, got %d", i, expectedBlock, gotBlock)
		}
		// got not nil, expected nil
		if gotBlock != nil && expectedBlock == nil {
			t.Fatalf("expected %d, got %d", expectedBlock, gotBlock)
		}
		// got not nil, expected not nil
		if expectedBlock.data != gotBlock.data {
			t.Fatalf("expected %d, got %d", expectedBlock, gotBlock)
		}
	}
}

func TestCompactByFile(t *testing.T) {
	disk := &Disk{
		blocks: []*Block{
			{data: 0}, {data: 0},
			nil, nil, nil,
			{data: 1}, {data: 1}, {data: 1},
			nil, nil, nil,
			{data: 2},
			nil, nil, nil,
			{data: 3}, {data: 3}, {data: 3},
			nil,
			{data: 4}, {data: 4},
			nil,
			{data: 5}, {data: 5}, {data: 5}, {data: 5},
			nil,
			{data: 6}, {data: 6}, {data: 6}, {data: 6},
			nil,
			{data: 7}, {data: 7}, {data: 7},
			nil,
			{data: 8}, {data: 8}, {data: 8}, {data: 8},
			{data: 9}, {data: 9},
		},
	}

	CompactByFile(disk)

	expected := []*Block{
		{data: 0}, {data: 0},
		{data: 9}, {data: 9},
		{data: 2},
		{data: 1}, {data: 1}, {data: 1},
		{data: 7}, {data: 7}, {data: 7},
		nil,
		{data: 4}, {data: 4},
		nil,
		{data: 3}, {data: 3}, {data: 3},
		nil, nil, nil, nil,
		{data: 5}, {data: 5}, {data: 5}, {data: 5},
		nil,
		{data: 6}, {data: 6}, {data: 6}, {data: 6},
		nil, nil, nil, nil, nil,
		{data: 8}, {data: 8}, {data: 8}, {data: 8},
		nil, nil,
	}

	for i := range expected {
		gotBlock := disk.blocks[i]
		expectedBlock := expected[i]

		// both nil
		if gotBlock == nil && expectedBlock == nil {
			continue
		}
		// got nil, expected not nil
		if gotBlock == nil && expectedBlock != nil {
			t.Fatalf("block[%d]: expected %d, got %d", i, expectedBlock, gotBlock)
		}
		// got not nil, expected nil
		if gotBlock != nil && expectedBlock == nil {
			t.Fatalf("expected %d, got %d", expectedBlock, gotBlock)
		}
		// got not nil, expected not nil
		if expectedBlock.data != gotBlock.data {
			t.Fatalf("expected %d, got %d", expectedBlock, gotBlock)
		}
	}
}

func TestChecksum(t *testing.T) {
	disk := &Disk{
		blocks: []*Block{
			{data: 0},
			{data: 0},
			{data: 9},
			{data: 9},
			{data: 8},
			{data: 1},
			{data: 1},
			{data: 1},
			{data: 8},
			{data: 8},
			{data: 8},
			{data: 2},
			{data: 7},
			{data: 7},
			{data: 7},
			{data: 3},
			{data: 3},
			{data: 3},
			{data: 6},
			{data: 4},
			{data: 4},
			{data: 6},
			{data: 5},
			{data: 5},
			{data: 5},
			{data: 5},
			{data: 6},
			{data: 6},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}

	expected := 1928
	actual := disk.checksum()
	if actual != expected {
		t.Errorf("expected checksum %d, got %d", expected, actual)
	}
}
