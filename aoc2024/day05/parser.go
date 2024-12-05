package main

import (
	"fmt"
	"strconv"
)

type PrecedenceRule struct {
	Before int
	After  int
}

type Page struct {
	Number int
}

type PrintJob struct {
	Pages []Page
}

type Program struct {
	Rules     []PrecedenceRule
	PrintJobs []PrintJob
}

type Parser struct {
	tokens []Token
	cur    int
}

func Parse(tokens []Token) (Program, error) {
	parser := Parser{tokens: tokens, cur: 0}
	return parser.parse()
}

func (p *Parser) parse() (Program, error) {
	program := Program{}
	for p.cur < len(p.tokens) {
		rule, err := p.parseRule()
		if err != nil {
			return program, err
		}
		program.Rules = append(program.Rules, rule)

		// a newline after a rule means the rules are done
		if p.peek().Type == "newline" {
			p.next()
			break
		}
	}

	for p.cur < len(p.tokens) {
		job, err := p.parsePrintJob()
		if err != nil {
			return program, err
		}
		program.PrintJobs = append(program.PrintJobs, job)
	}

	return program, nil
}

func (p *Parser) parseRule() (PrecedenceRule, error) {
	before, err := p.parseInt()
	if err != nil {
		return PrecedenceRule{}, err
	}

	if err := p.expect("pipe"); err != nil {
		return PrecedenceRule{}, err
	}

	after, err := p.parseInt()
	if err != nil {
		return PrecedenceRule{}, err
	}

	if err := p.expect("newline"); err != nil {
		return PrecedenceRule{}, err
	}

	return PrecedenceRule{Before: before, After: after}, nil
}

func (p *Parser) parsePrintJob() (PrintJob, error) {
	var pages []Page
loop:
	for p.cur < len(p.tokens) {
		page, err := p.parsePage()
		if err != nil {
			return PrintJob{}, err
		}
		pages = append(pages, page)

		switch p.peek().Type {
		case "newline":
			p.next()
			break loop
		case "comma":
			p.next()
			continue loop
		default:
			if p.cur >= len(p.tokens) {
				break loop
			}
			return PrintJob{}, fmt.Errorf("expected newline or comma, got %s", p.peek().Type)
		}
	}
	return PrintJob{Pages: pages}, nil
}

func (p *Parser) parsePage() (Page, error) {
	number, err := p.parseInt()
	if err != nil {
		return Page{}, err
	}
	return Page{Number: number}, nil
}

func (p *Parser) expect(tokenType string) error {
	token := p.peek()
	if token.Type != tokenType {
		return fmt.Errorf("expected %s, got %s (cur: %d)", tokenType, token.Type, p.cur)
	}
	p.next()
	return nil
}

func (p *Parser) parseInt() (int, error) {
	token := p.peek()
	if token.Type != "int" {
		fmt.Println(token, p.cur)
		return 0, fmt.Errorf("expected int, got %s", token.Type)
	}
	value, err := strconv.Atoi(string(token.Value))
	if err != nil {
		return 0, err
	}
	p.next()
	return value, nil
}

func (p *Parser) peek() Token {
	if p.cur >= len(p.tokens) {
		return Token{}
	}
	return p.tokens[p.cur]
}

func (p *Parser) peek2() (Token, Token) {
	if p.cur+1 >= len(p.tokens) {
		return Token{}, Token{}
	}
	return p.tokens[p.cur], p.tokens[p.cur+1]
}

func (p *Parser) next() Token {
	if p.cur >= len(p.tokens) {
		return Token{}
	}
	t := p.tokens[p.cur]
	p.cur++
	return t
}
