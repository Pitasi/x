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

	for p.cur < len(p.input) {
		row := p.parseRow()
		rows = append(rows, row)
	}

	return Map{
		grid: rows,
	}
}

func (p *parser) parseRow() []Cell {
	var cells []Cell

	var x int
	for p.cur < len(p.input) && p.input[p.cur] != '\n' {
		switch p.input[p.cur] {
		case '.':
			cells = append(cells, Cell{})
		default:
			cells = append(cells, Cell{antenna: p.input[p.cur]})
		}
		p.cur++
		x++
	}

	p.cur++ // skip newline

	return cells
}
