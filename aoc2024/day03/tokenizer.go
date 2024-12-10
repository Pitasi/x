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

	mulDisabled bool
}

type Mul struct {
	A, B int
}

func (t *Tokenizer) Next() (Mul, bool) {
	for t.cur < len(t.input) {
		if t.mulDisabled {
			if found := t.consumeNonD(); !found {
				break
			}
			if t.tryConsumeDo() {
				t.mulDisabled = false
			}
		} else {
			if found := t.consumeNonMD(); !found {
				break
			}
			if t.input[t.cur] == 'd' {
				if t.tryConsumeDont() {
					t.mulDisabled = true
					continue
				}
			} else if t.input[t.cur] == 'm' {
				mul, ok := t.tryConsumeMul()
				if ok {
					return mul, true
				}
			}
		}
	}
	return Mul{}, false
}

func (t *Tokenizer) consumeNonD() bool {
	for t.cur < len(t.input) && t.input[t.cur] != 'd' {
		t.cur++
	}
	return t.cur < len(t.input)
}

func (t *Tokenizer) consumeNonMD() bool {
	for t.cur < len(t.input) && t.input[t.cur] != 'm' && t.input[t.cur] != 'd' {
		t.cur++
	}
	return t.cur < len(t.input)
}

func (t *Tokenizer) tryConsumeMul() (Mul, bool) {
	if !t.expectChar('m') {
		return Mul{}, false
	}
	if !t.expectChar('u') {
		return Mul{}, false
	}
	if !t.expectChar('l') {
		return Mul{}, false
	}
	if !t.expectChar('(') {
		return Mul{}, false
	}
	a, ok := t.tryConsumeInt()
	if !ok {
		return Mul{}, false
	}
	if !t.expectChar(',') {
		return Mul{}, false
	}
	b, ok := t.tryConsumeInt()
	if !ok {
		return Mul{}, false
	}
	if !t.expectChar(')') {
		return Mul{}, false
	}
	return Mul{a, b}, true
}

func (t *Tokenizer) tryConsumeDont() bool {
	if !t.expectChar('d') {
		return false
	}
	if !t.expectChar('o') {
		return false
	}
	if !t.expectChar('n') {
		return false
	}
	if !t.expectChar('\'') {
		return false
	}
	if !t.expectChar('t') {
		return false
	}
	if !t.expectChar('(') {
		return false
	}
	if !t.expectChar(')') {
		return false
	}
	return true
}

func (t *Tokenizer) tryConsumeDo() bool {
	if !t.expectChar('d') {
		return false
	}
	if !t.expectChar('o') {
		return false
	}
	if !t.expectChar('(') {
		return false
	}
	if !t.expectChar(')') {
		return false
	}
	return true
}

func (t *Tokenizer) tryConsumeInt() (int, bool) {
	start := t.cur
	for t.cur < len(t.input) && t.input[t.cur] >= '0' && t.input[t.cur] <= '9' {
		t.cur++
	}
	if t.cur == start {
		return 0, false
	}

	n, err := strconv.Atoi(string(t.input[start:t.cur]))
	if err != nil {
		panic(err)
	}

	return n, true
}

func (t *Tokenizer) expectChar(c byte) bool {
	if t.cur >= len(t.input) {
		return false
	}
	ok := t.input[t.cur] == c
	t.cur++
	return ok
}
