package main

import (
	"fmt"
	"iter"
	"strconv"
)

type Equation struct {
	Result int
	Inputs []int
}

func parse(input []byte) iter.Seq[Equation] {
	p := Parser{input, 0}
	return p.parse()
}

type Parser struct {
	input []byte
	cur   int
}

func (p *Parser) parse() iter.Seq[Equation] {
	return func(yield func(Equation) bool) {
		for p.cur < len(p.input) {
			if !p.parseEquation(yield) {
				break
			}
		}
	}
}

func (p *Parser) parseEquation(yield func(Equation) bool) bool {
	result, err := p.parseInt()
	if err != nil {
		return false
	}

	if !p.expect(':') {
		return false
	}

	inputs, err := p.parseInts()
	if err != nil {
		return false
	}

	yield(Equation{result, inputs})
	return true
}

func (p *Parser) parseInt() (int, error) {
	start := p.cur
	for p.cur < len(p.input) && p.input[p.cur] >= '0' && p.input[p.cur] <= '9' {
		p.cur++
	}

	if start == p.cur {
		return 0, fmt.Errorf("expected int, got %s", p.input[start:p.cur])
	}

	return strconv.Atoi(string(p.input[start:p.cur]))
}

func (p *Parser) expect(b byte) bool {
	if p.cur >= len(p.input) {
		return false
	}
	if p.input[p.cur] != b {
		return false
	}
	p.cur++
	return true
}

func (p *Parser) parseInts() ([]int, error) {
	var inputs []int
	for p.cur < len(p.input) {
		if !p.expect(' ') {
			return nil, fmt.Errorf("expected space, got %s", p.input[p.cur:p.cur+1])
		}

		input, err := p.parseInt()
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)

		if p.expect('\n') {
			break
		}
	}
	return inputs, nil
}
