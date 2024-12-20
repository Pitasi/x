package main

import (
	"fmt"
	"io"
	"testing"
)

func TestPathExists(t *testing.T) {
	input := []byte(`r, wr, b, g, bwu, rb, gb, br

brwrr`)

	available, _ := parse(input)

	tr := newTrie()
	for _, d := range available {
		tr.add(d)
	}

	if PathExists(tr, Design("brwrr")) == 0 {
		t.Fatal("brwrr should exist")
	}

	if PathExists(tr, Design("ubwu")) > 0 {
		t.Fatal("ubwu should not exist")
	}

	if PathExists(tr, Design("bwurrg")) == 0 {
		t.Fatal("bwurrg should exist")
	}
}

func TestExample(t *testing.T) {
	input := []byte(`r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`)

	available, requested := parse(input)
	tr := newTrie()
	for _, d := range available {
		tr.add(d)
	}

	var count int
	for _, d := range requested {
		count += PathExists(tr, d)
	}

	expected := 16
	if count != expected {
		t.Fatalf("expected %d got %d", expected, count)
	}
}

func BenchmarkFull(b *testing.B) {
	for range b.N {
		available, requested := parse(input)
		tr := newTrie()
		for _, d := range available {
			tr.add(d)
		}

		var count int
		for _, d := range requested {
			count += PathExists(tr, d)
		}
		fmt.Fprint(io.Discard, count)
	}
}
