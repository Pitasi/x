package main

type Vector struct {
	X int
	Y int
}

func (v Vector) Add(other Vector) Vector {
	return Vector{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vector) Mul(n int) Vector {
	return Vector{X: v.X * n, Y: v.Y * n}
}

func (v Vector) Mod(bounds Vector) Vector {
	v2 := Vector{
		X: v.X % bounds.X,
		Y: v.Y % bounds.Y,
	}
	if v2.X < 0 {
		v2.X += bounds.X
	}
	if v2.Y < 0 {
		v2.Y += bounds.Y
	}
	return v2
}

func v(x, y int) Vector {
	return Vector{X: x, Y: y}
}

type Robot struct {
	P Vector
	V Vector
}

func Parse(in []byte) []Robot {
	p := parser{in: in}
	return p.parse()
}

type parser struct {
	in  []byte
	cur int
}

func (p *parser) parse() []Robot {
	robots := make([]Robot, 0)
	for p.cur < len(p.in) {
		robots = append(robots, p.parseRobot())
	}
	return robots
}

func (p *parser) parseRobot() Robot {
	p.expect('p')
	p.expect('=')
	x, y := p.parseVector()
	p.expect(' ')
	p.expect('v')
	p.expect('=')
	vx, vy := p.parseVector()
	p.newlineOrEOF()
	return Robot{P: v(x, y), V: v(vx, vy)}
}

func (p *parser) parseVector() (int, int) {
	x := p.parseInt()
	p.expect(',')
	y := p.parseInt()
	return x, y
}

func (p *parser) parseInt() int {
	var v int

	sign := 1
	switch p.in[p.cur] {
	case '-':
		sign = -1
		p.cur++
	case '+':
		p.cur++
	}

	for i := p.cur; i < len(p.in); i++ {
		c := p.in[i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int(c-'0')
		p.cur++
	}

	return sign * v
}

func (p *parser) expect(c byte) {
	if p.cur >= len(p.in) {
		panic("unexpected EOF")
	}
	if p.in[p.cur] != c {
		panic("expected '" + string(c) + "' got '" + string(p.in[p.cur]) + "'")
	}
	p.cur++
}

func (p *parser) newlineOrEOF() {
	if p.cur >= len(p.in) {
		return
	}
	if p.in[p.cur] == '\n' {
		p.cur++
	}
}
