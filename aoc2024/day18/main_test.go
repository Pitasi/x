package main

import (
	"testing"
)

func TestExample(t *testing.T) {
	input := []byte(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`)

	g, blocks := Build(7, input)
	blocks = blocks[:12]
	for _, b := range blocks {
		g[b.y][b.x] = '#'
	}
	path := astar2(g, 6, 6)
	expected := 22
	if len(path)-1 != expected {
		t.Errorf("expected %d, got %d", expected, len(path)-1)
	}
}

func TestExample2(t *testing.T) {
	input := []byte(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`)

	g, blocks := Build(7, input)
	var block vec
	for _, b := range blocks {
		g[b.y][b.x] = '#'
		path := astar2(g, 6, 6)
		if len(path) == 0 {
			block = b
			break
		}
	}
	expected := vec{x: 6, y: 1}
	if block != expected {
		t.Errorf("expected %d, got %d", expected, block)
	}
}
