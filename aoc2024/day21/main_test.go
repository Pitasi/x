package main

import (
	"testing"
)

func TestStar1(t *testing.T) {
	cases := []struct {
		input []byte
		score int
	}{
		{[]byte("029A"), 68 * 29},
		{[]byte("980A"), 60 * 980},
		{[]byte("179A"), 68 * 179},
		{[]byte("456A"), 64 * 456},
		{[]byte("379A"), 64 * 379},
	}

	for _, c := range cases {
		t.Run(string(c.input), func(t *testing.T) {
			score := day21(c.input, 1)
			if score != c.score {
				t.Errorf("expected %d, got %d", c.score, score)
			}
		})
	}
}
