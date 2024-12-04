package main

import "testing"

func TestCountXmas(t *testing.T) {
	g := Grid{
		rows: [][]byte{
			{'M', 'M', 'M', 'S', 'X', 'X', 'M', 'A', 'S', 'M'},
			{'M', 'S', 'A', 'M', 'X', 'M', 'S', 'M', 'S', 'A'},
			{'A', 'M', 'X', 'S', 'X', 'M', 'A', 'A', 'M', 'M'},
			{'M', 'S', 'A', 'M', 'A', 'S', 'M', 'S', 'M', 'X'},
			{'X', 'M', 'A', 'S', 'A', 'M', 'X', 'A', 'M', 'M'},
			{'X', 'X', 'A', 'M', 'M', 'X', 'X', 'A', 'M', 'A'},
			{'S', 'M', 'S', 'M', 'S', 'A', 'S', 'X', 'S', 'S'},
			{'S', 'A', 'X', 'A', 'M', 'A', 'S', 'A', 'A', 'A'},
			{'M', 'A', 'M', 'M', 'M', 'X', 'M', 'M', 'M', 'M'},
			{'M', 'X', 'M', 'X', 'A', 'X', 'M', 'A', 'S', 'X'},
		},
	}
	count := CountXmas(g)
	if count != 18 {
		t.Errorf("CountXmas() = %v, want %v", count, 18)
	}
}

func TestCountXs(t *testing.T) {
	g := Grid{
		rows: [][]byte{
			{'M', 'M', 'M', 'S', 'X', 'X', 'M', 'A', 'S', 'M'},
			{'M', 'S', 'A', 'M', 'X', 'M', 'S', 'M', 'S', 'A'},
			{'A', 'M', 'X', 'S', 'X', 'M', 'A', 'A', 'M', 'M'},
			{'M', 'S', 'A', 'M', 'A', 'S', 'M', 'S', 'M', 'X'},
			{'X', 'M', 'A', 'S', 'A', 'M', 'X', 'A', 'M', 'M'},
			{'X', 'X', 'A', 'M', 'M', 'X', 'X', 'A', 'M', 'A'},
			{'S', 'M', 'S', 'M', 'S', 'A', 'S', 'X', 'S', 'S'},
			{'S', 'A', 'X', 'A', 'M', 'A', 'S', 'A', 'A', 'A'},
			{'M', 'A', 'M', 'M', 'M', 'X', 'M', 'M', 'M', 'M'},
			{'M', 'X', 'M', 'X', 'A', 'X', 'M', 'A', 'S', 'X'},
		},
	}
	count := CountXs(g)
	if count != 9 {
		t.Errorf("CountXs() = %v, want %v", count, 9)
	}
}
