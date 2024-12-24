package main

func Parse(input []byte) map[string]Gate {
	p := parser{input: input}
	return p.parse()
}

type parser struct {
	input []byte
	pos   int
}

func (p *parser) parse() map[string]Gate {
	g := make(map[string]Gate)
	for p.pos < len(p.input) {
		name := p.parseName()
		if p.input[p.pos] != ':' {
			panic("expected :")
		}
		p.pos++
		if p.input[p.pos] != ' ' {
			panic("expected space")
		}
		p.pos++

		g[name] = p.parseBool()

		if p.input[p.pos] != '\n' {
			panic("expected newline")
		}
		p.pos++
		if p.pos < len(p.input) && p.input[p.pos] == '\n' {
			p.pos++
			break
		}
	}

	for p.pos < len(p.input) {
		nameLeft := p.parseName()
		if p.input[p.pos] != ' ' {
			panic("expected space")
		}
		p.pos++

		op := p.parseOp()
		if p.input[p.pos] != ' ' {
			panic("expected space")
		}
		p.pos++

		nameRight := p.parseName()
		if p.input[p.pos] != ' ' {
			panic("expected space")
		}
		p.pos++

		if p.input[p.pos] != '-' {
			panic("expected -")
		}
		p.pos++

		if p.input[p.pos] != '>' {
			panic("expected >")
		}
		p.pos++

		if p.input[p.pos] != ' ' {
			panic("expected space")
		}
		p.pos++

		name := p.parseName()
		if p.pos < len(p.input) && p.input[p.pos] == '\n' {
			p.pos++
		}

		gLeft := &LazyGate{m: g, name: nameLeft}

		gRight := &LazyGate{m: g, name: nameRight}

		switch op {
		case "AND":
			g[name] = &GateAnd{left: gLeft, right: gRight}
		case "OR":
			g[name] = &GateOr{left: gLeft, right: gRight}
		case "XOR":
			g[name] = &GateXor{left: gLeft, right: gRight}
		default:
			panic("invalid op")
		}
	}

	return g
}

func (p *parser) parseName() string {
	start := p.pos
	end := p.pos + 3
	name := p.input[start:end]
	p.pos += 3
	return string(name)
}

func (p *parser) parseBool() Gate {
	if p.input[p.pos] == '1' {
		p.pos++
		return &GateTrue{}
	}
	if p.input[p.pos] == '0' {
		p.pos++
		return &GateFalse{}
	}

	panic("expected 1 or 0")
}

func (p *parser) parseOp() string {
	peek := string(p.input[p.pos : p.pos+3])
	switch peek {
	case "AND":
		p.pos += 3
		return "AND"
	case "OR ":
		p.pos += 2
		return "OR"
	case "XOR":
		p.pos += 3
		return "XOR"
	default:
		panic("expected AND, OR, or XOR")
	}
}
