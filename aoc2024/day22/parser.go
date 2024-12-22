package main

import "strconv"

func Parse(in []byte) []int {
	p := parser{input: in}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

func (p *parser) parse() []int {
	var result []int

	for p.pos < len(p.input) {
		result = append(result, p.parseInt())
	}

	return result
}

func (p *parser) parseInt() int {
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] >= '0' && p.input[p.pos] <= '9' {
		p.pos++
	}

	result, err := strconv.Atoi(string(p.input[start:p.pos]))
	if err != nil {
		panic(err)
	}

	if p.pos < len(p.input) && p.input[p.pos] == '\n' {
		p.pos++
	}

	return result
}
