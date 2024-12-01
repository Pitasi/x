package main

import (
	_ "embed"
	"strconv"
)

//go:embed input.txt
var input []byte

type Tokenizer struct {
	input []byte
	cur   int
}

func (t *Tokenizer) Next() (int, int, bool) {
	if t.cur >= len(t.input) {
		return 0, 0, false
	}
	left := t.consumeInt()
	t.consumeWhitespaces()
	right := t.consumeInt()
	t.consumeNewline()
	return left, right, true
}

func (t *Tokenizer) consumeInt() int {
	start := t.cur
	for t.cur < len(t.input) && isDigit(t.input[t.cur]) {
		t.cur++
	}

	i, err := strconv.Atoi(string(t.input[start:t.cur]))
	if err != nil {
		panic(err)
	}

	return i
}

func (t *Tokenizer) consumeWhitespaces() {
	for t.cur < len(t.input) && isWhitespace(t.input[t.cur]) {
		t.cur++
	}
}

func (t *Tokenizer) consumeNewline() {
	if t.input[t.cur] == '\n' {
		t.cur++
	}
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
