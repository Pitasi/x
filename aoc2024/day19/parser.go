package main

type Design []byte

func (d Design) String() string {
	return string(d)
}

func parse(input []byte) ([]Design, []Design) {
	p := parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

func (p *parser) parse() ([]Design, []Design) {
	var (
		available []Design
		requested []Design
	)
	for p.pos < len(p.input) {
		design := p.parseDesign()
		available = append(available, design)
		if p.input[p.pos] == '\n' {
			break
		}
		p.expect(", ")
	}

	if p.input[p.pos] == '\n' {
		p.pos++
	} else {
		panic("missing newline")
	}
	if p.input[p.pos] == '\n' {
		p.pos++
	} else {
		panic("missing newline")
	}

	for p.pos < len(p.input) {
		design := p.parseDesign()
		requested = append(requested, design)
		if p.pos < len(p.input) && p.input[p.pos] == '\n' {
			p.pos++
		}
	}

	return available, requested
}

func (p *parser) parseDesign() Design {
	start := p.pos
	for p.pos < len(p.input) && isColor(byte(p.input[p.pos])) {
		p.pos++
	}
	return Design(p.input[start:p.pos])
}

func isColor(r byte) bool {
	return r == 'w' || r == 'u' || r == 'b' || r == 'g' || r == 'r'
}

func (p *parser) expect(s string) {
	for _, c := range s {
		if p.input[p.pos] != byte(c) {
			panic("expected " + string(c))
		}
		p.pos++
	}
}
