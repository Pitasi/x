package main

func Parse(input []byte) [][]*Node {
	p := &parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
}

func (p *parser) parse() [][]*Node {
	grid, width, height := p.parseGrid()

	nodesGrid := make([][]*Node, height)
	for i := range width {
		nodesGrid[i] = make([]*Node, width)
	}

	for y := range height {
		for x := range width {
			nodesGrid[y][x] = &Node{
				Value: grid[y][x],
				X:     x,
				Y:     y,
			}
		}
	}

	for x := range width {
		for y := range height {
			if y > 0 {
				nodesGrid[y][x].Up = nodesGrid[y-1][x]
			}
			if y < height-1 {
				nodesGrid[y][x].Down = nodesGrid[y+1][x]
			}
			if x > 0 {
				nodesGrid[y][x].Left = nodesGrid[y][x-1]
			}
			if x < width-1 {
				nodesGrid[y][x].Right = nodesGrid[y][x+1]
			}
		}
	}

	return nodesGrid
}

func (p *parser) parseGrid() ([][]byte, int, int) {
	var cur int
	var (
		rows   [][]byte
		curRow []byte
	)
	for cur < len(p.input) {
		if p.input[cur] == '\n' {
			rows = append(rows, curRow)
			curRow = []byte{}
		} else {
			curRow = append(curRow, p.input[cur])
		}
		cur++
	}
	if len(curRow) > 0 {
		rows = append(rows, curRow)
	}
	height := len(rows)
	width := len(rows[0])
	return rows, width, height
}
