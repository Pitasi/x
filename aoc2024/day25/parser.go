package main

import "bytes"

type Key struct {
	a, b, c, d, e int
	isLock        bool
}

func (k Key) Fit(other Key) bool {
	return k.isLock != other.isLock && k.a+other.a <= 5 && k.b+other.b <= 5 && k.c+other.c <= 5 && k.d+other.d <= 5 && k.e+other.e <= 5
}

func Parse(input []byte) []Key {
	p := parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

func (p *parser) parse() []Key {
	var keys []Key

	for p.pos < len(p.input) {
		keys = append(keys, p.parseKey())
		if p.pos < len(p.input) && p.input[p.pos] == '\n' {
			p.pos++
		}
	}

	return keys
}

func (p *parser) parseKey() Key {
	firstRow := p.parseRow()
	isLock := bytes.Equal(firstRow, []byte("#####"))
	var rows [][]byte
	for i := 0; i < 5; i++ {
		rows = append(rows, p.parseRow())
	}
	p.parseRow()
	a, b, c, d, e := p.countColumns(rows)
	return Key{a, b, c, d, e, isLock}
}

func (p *parser) parseRow() []byte {
	start := p.pos
	length := 5
	row := p.input[start : start+length]
	p.pos += length
	if p.pos < len(p.input) && p.input[p.pos] == '\n' {
		p.pos++
	}
	return row
}

func (p *parser) countColumns(rows [][]byte) (int, int, int, int, int) {
	counts := make([]int, 5)
	for y := range rows {
		for x := range rows[y] {
			if rows[y][x] == '#' {
				counts[x]++
			}
		}
	}
	return counts[0], counts[1], counts[2], counts[3], counts[4]
}
