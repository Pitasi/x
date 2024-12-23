package main

func Parse(input []byte) Graph {
	p := Parser{input: input}
	return p.parse()
}

type Parser struct {
	input []byte
	pos   int
}

func (p *Parser) parse() Graph {
	g := newGraph()
	for p.pos < len(p.input) {
		node1, node2 := p.parseLine()
		g.AddNode(node1)
		g.AddNode(node2)
		g.AddEdge(node1, node2)
	}
	return g
}

func (p *Parser) parseLine() (string, string) {
	node1 := p.input[p.pos : p.pos+2]
	p.pos += 2
	if p.input[p.pos] != '-' {
		panic("expected -")
	}
	p.pos++
	node2 := p.input[p.pos : p.pos+2]
	p.pos += 2
	if p.pos < len(p.input) && p.input[p.pos] == '\n' {
		p.pos++
	}
	return string(node1), string(node2)
}
