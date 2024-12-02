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

func (t *Tokenizer) Next() ([]int, bool) {
	if t.cur >= len(t.input) {
		return nil, false
	}
	ints := t.consumeInts()
	return ints, true
}

func (t *Tokenizer) consumeInts() []int {
	ints := make([]int, 0)

	for {
		start := t.cur
		for t.cur < len(t.input) && isDigit(t.input[t.cur]) {
			t.cur++
		}

		i, err := strconv.Atoi(string(t.input[start:t.cur]))
		if err != nil {
			panic(err)
		}

		ints = append(ints, i)

		t.consumeWhitespaces()
		if t.consumeNewline() || t.cur >= len(t.input) {
			break
		}
	}

	return ints
}

func (t *Tokenizer) consumeWhitespaces() {
	for t.cur < len(t.input) && isWhitespace(t.input[t.cur]) {
		t.cur++
	}
}

func (t *Tokenizer) consumeNewline() bool {
	if t.input[t.cur] == '\n' {
		t.cur++
		return true
	}
	return false
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
