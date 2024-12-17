package main

import (
	"fmt"
	"strconv"
)

func Parse(input []byte) (int, int, int, []byte) {
	p := parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

func (p *parser) parse() (int, int, int, []byte) {
	a := p.parseReg('A')
	b := p.parseReg('B')
	c := p.parseReg('C')
	p.consumeNewline()
	program := p.parseProg()
	return a, b, c, program
}

func (p *parser) parseReg(r rune) int {
	p.expect("Register ")
	p.expect(string(r))
	p.expect(": ")
	v := p.parseInt()
	p.consumeNewline()
	return v
}

func (p *parser) parseProg() []byte {
	p.expect("Program: ")
	var ops []byte
	for p.pos < len(p.input) {
		ops = append(ops, p.parseOp())
		p.consumeNewline()
	}
	return ops
}

func (p *parser) parseOp() byte {
	if p.input[p.pos] < '0' || p.input[p.pos] > '9' {
		panic(fmt.Sprintf("invalid opcode at pos %d/%d: '%c'", p.pos, len(p.input), p.input[p.pos]))
	}

	op := p.input[p.pos] - '0'
	p.pos++

	if p.pos < len(p.input)-1 {
		p.expect(",")
	}

	return op
}

func (p *parser) parseInt() int {
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] >= '0' && p.input[p.pos] <= '9' {
		p.pos++
	}
	end := p.pos
	s := string(p.input[start:end])
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("invalid integer")
	}
	return i
}

func (p *parser) expect(s string) {
	if p.pos >= len(p.input) {
		panic("unexpected end of input")
	}
	if string(p.input[p.pos:p.pos+len(s)]) != s {
		panic(fmt.Sprintf("expected %q, got %q", s, p.input[p.pos:p.pos+len(s)]))
	}
	p.pos += len(s)
}

func (p *parser) consumeNewline() {
	if p.pos >= len(p.input) {
		return
	}
	if p.input[p.pos] == '\n' {
		p.pos++
	}
}
