package main

func Parse(input []byte) Map {
	p := parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	cur   int
}

func (p *parser) parse() Map {
	var rows [][]Cell
	var guard *Guard

	var y int
	for p.cur < len(p.input) {
		row, maybeGuard := p.parseRow(y)
		rows = append(rows, row)
		if maybeGuard != nil {
			guard = maybeGuard
		}
		y++
	}

	if guard == nil {
		panic("guard not found")
	}

	return Map{
		grid:  rows,
		guard: *guard,
	}
}

func (p *parser) parseRow(y int) ([]Cell, *Guard) {
	var (
		cells []Cell
		guard *Guard
	)

	var x int
	for p.cur < len(p.input) && p.input[p.cur] != '\n' {
		switch p.input[p.cur] {
		case '.':
			cells = append(cells, &Empty{})
		case '#':
			cells = append(cells, &Wall{})
		case '^':
			guard = &Guard{
				position: Position{
					x: x,
					y: y,
				},
				direction: Up,
			}
			cells = append(cells, &Empty{})
		}
		p.cur++
		x++
	}

	p.cur++ // skip newline

	return cells, guard
}
