package main

import (
	"fmt"
	"strings"
)

func Parse(input []byte) (Map, Position, []Move) {
	p := &parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

type Position struct {
	X int
	Y int
}

type Map [][]byte

func (m Map) String() string {
	var sb strings.Builder
	for _, row := range m {
		for _, c := range row {
			sb.WriteByte(c)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type Move byte

func (m Move) String() string {
	switch m {
	case Up:
		return "^"
	case Down:
		return "v"
	case Left:
		return "<"
	case Right:
		return ">"
	}
	return fmt.Sprintf("ERR{%v}", byte(m))
}

const (
	Up    Move = '^'
	Down  Move = 'v'
	Left  Move = '<'
	Right Move = '>'
)

func (p *parser) parse() (Map, Position, []Move) {
	m, robot, err := p.parseMap()
	if err != nil {
		panic(err)
	}

	moves, err := p.parseMoves()
	if err != nil {
		panic(err)
	}

	return m, robot, moves
}

func (p *parser) parseMap() (Map, Position, error) {
	var (
		m          Map
		currentRow []byte
		robot      Position
	)

	for p.pos < len(p.input) {
		c := p.input[p.pos]
		if c == '@' {
			robot.Y = len(m)
			robot.X = len(currentRow)
		}
		if c == '\n' {
			if len(currentRow) == 0 {
				break
			}
			m = append(m, currentRow)
			currentRow = nil
		} else {
			currentRow = append(currentRow, c)
		}
		p.pos++
	}

	if p.input[p.pos] != '\n' {
		return nil, Position{}, fmt.Errorf("expected newline at the end of the input")
	}
	p.pos++

	return m, robot, nil
}

func (p *parser) parseMoves() ([]Move, error) {
	var moves []Move

	for p.pos < len(p.input) {
		var move Move
		switch p.input[p.pos] {
		case '^':
			move = Up
		case 'v':
			move = Down
		case '<':
			move = Left
		case '>':
			move = Right
		case '\n':
			p.pos++
			continue
		default:
			return nil, fmt.Errorf("unexpected character '%c'", p.input[p.pos])
		}
		moves = append(moves, move)
		p.pos++
	}

	return moves, nil
}
