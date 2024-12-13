package main

import "fmt"

type Position struct {
	X int
	Y int
}

func (p Position) Add(q Position) Position {
	return Position{X: p.X + q.X, Y: p.Y + q.Y}
}

type Game struct {
	Prize Position

	A Position
	B Position
}

func Parse(input []byte) ([]Game, error) {
	p := parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

func (p *parser) parse() ([]Game, error) {
	var games []Game
	for p.pos < len(p.input) {
		g, err := p.parseGame()
		if err != nil {
			return nil, fmt.Errorf("parsing game at position %d: %w", p.pos, err)
		}
		games = append(games, g)
	}
	return games, nil
}

func (p *parser) parseGame() (Game, error) {
	var game Game
	btnA, err := p.parseButton('A')
	if err != nil {
		return game, err
	}
	game.A = btnA
	btnB, err := p.parseButton('B')
	if err != nil {
		return game, err
	}
	game.B = btnB

	prize, err := p.parsePrize()
	if err != nil {
		return game, err
	}
	game.Prize = prize

	if p.pos < len(p.input) && p.input[p.pos] == '\n' {
		p.pos++
	}
	if p.pos < len(p.input) && p.input[p.pos] == '\n' {
		p.pos++
	}

	return game, nil
}

func (p *parser) parseButton(btn rune) (Position, error) {
	if err := p.expect("Button "); err != nil {
		return Position{}, err
	}
	if err := p.expect(string(btn)); err != nil {
		return Position{}, err
	}
	if err := p.expect(": "); err != nil {
		return Position{}, err
	}

	x, err := p.parseX()
	if err != nil {
		return Position{}, err
	}

	if err := p.expect(", "); err != nil {
		return Position{}, err
	}

	y, err := p.parseY()
	if err != nil {
		panic(err)
	}

	if p.pos < len(p.input) && p.input[p.pos] == byte('\n') {
		p.pos++
	}

	pos := Position{X: x, Y: y}
	return pos, nil
}

func (p *parser) parseX() (int, error) {
	return p.parseDelta('X')
}

func (p *parser) parseY() (int, error) {
	return p.parseDelta('Y')
}

func (p *parser) parseDelta(c byte) (int, error) {
	if err := p.expect(string(c)); err != nil {
		return 0, err
	}
	v, err := p.parseSignedInt()
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (p *parser) parsePrize() (Position, error) {
	if err := p.expect("Prize: "); err != nil {
		return Position{}, err
	}

	if err := p.expect("X="); err != nil {
		return Position{}, err
	}
	x, err := p.parseInt()
	if err != nil {
		return Position{}, err
	}
	if err := p.expect(", Y="); err != nil {
		return Position{}, err
	}
	y, err := p.parseInt()
	if err != nil {
		return Position{}, err
	}

	prize := Position{X: x, Y: y}
	return prize, nil
}

func (p *parser) parseSignedInt() (int, error) {
	var v int

	var sign int
	if p.input[p.pos] == '-' {
		sign = -1
		p.pos++
	} else if p.input[p.pos] == '+' {
		sign = 1
		p.pos++
	} else {
		return 0, fmt.Errorf("expected '+' or '-', got '%s'", string(p.input[p.pos]))
	}

	v, err := p.parseInt()
	if err != nil {
		return 0, err
	}

	v = v * sign
	return v, nil
}

func (p *parser) parseInt() (int, error) {
	var v int
	for i := p.pos; i < len(p.input); i++ {
		c := p.input[i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int(c-'0')
		p.pos++
	}
	return v, nil
}

func (p *parser) expect(s string) error {
	for _, c := range s {
		if p.input[p.pos] != byte(c) {
			return fmt.Errorf("expected '%s', got '%s'", s, string(p.input[p.pos]))
		}
		p.pos++
	}

	return nil
}
