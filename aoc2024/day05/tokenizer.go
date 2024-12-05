package main

type Token struct {
	Type  string
	Value []byte
}

type Tokenizer struct {
	input []byte
	cur   int
}

func Tokenize(input []byte) ([]Token, error) {
	tokenizer := Tokenizer{input: input, cur: 0}
	return tokenizer.tokenize()
}

func (t *Tokenizer) tokenize() ([]Token, error) {
	var tokens []Token
	for t.cur < len(t.input) {
		token, err := t.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (t *Tokenizer) nextToken() (Token, error) {
	switch t.peek() {
	case '|':
		return t.parsePipe()
	case ',':
		return t.parseComma()
	case '\n':
		return t.parseNewline()
	default:
		return t.parseInt()
	}
}

func (t *Tokenizer) parsePipe() (Token, error) {
	t.next()
	return Token{Type: "pipe", Value: []byte("|")}, nil
}

func (t *Tokenizer) parseComma() (Token, error) {
	t.next()
	return Token{Type: "comma", Value: []byte(",")}, nil
}

func (t *Tokenizer) parseNewline() (Token, error) {
	t.next()
	return Token{Type: "newline", Value: []byte("\n")}, nil
}

func (t *Tokenizer) parseInt() (Token, error) {
	start := t.cur
	for t.cur < len(t.input) && t.isDigit(t.input[t.cur]) {
		t.cur++
	}
	value := t.input[start:t.cur]
	return Token{Type: "int", Value: value}, nil
}

func (t *Tokenizer) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (t *Tokenizer) peek() byte {
	if t.cur >= len(t.input) {
		return 0
	}
	return t.input[t.cur]
}

func (t *Tokenizer) next() byte {
	if t.cur >= len(t.input) {
		return 0
	}
	b := t.input[t.cur]
	t.cur++
	return b
}
