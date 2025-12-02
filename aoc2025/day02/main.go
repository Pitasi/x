package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	star(input, invalidID1)
	star(input, invalidID2)
}

func star(input []byte, checkfn func(string) bool) {
	s := &scanner{input: input}

	var count int
	for {
		// parse next range
		rangeStart := s.nextItem()
		expectType(rangeStart, itemNumber)
		dash := s.nextItem()
		expectType(dash, itemDash)
		rangeEnd := s.nextItem()
		expectType(rangeEnd, itemNumber)

		count += handleRange(rangeStart.s, rangeEnd.s, checkfn)

		commaOrEof := s.nextItem()
		if commaOrEof.typ == itemEOF {
			break
		}
		expectType(commaOrEof, itemComma)
	}

	fmt.Println("star:", count)
}

func handleRange(start, end string, checkfn func(string) bool) int {
	s, err := strconv.Atoi(start)
	if err != nil {
		panic(err)
	}
	e, err := strconv.Atoi(end)
	if err != nil {
		panic(err)
	}

	var count int
	for i := s; i <= e; i++ {
		id := strconv.Itoa(i)
		if checkfn(id) {
			count += i
		}
	}

	return count
}

func invalidID1(i string) bool {
	return i[:len(i)/2] == i[len(i)/2:]
}

func invalidID2(i string) bool {
	for length := 1; length < len(i); length++ {
		s := strings.Repeat(i[:length], len(i)/length)
		if s == i {
			return true
		}
	}
	return false
}

func expectType(it item, typ itemType) {
	if it.typ != typ {
		panic(fmt.Sprintf("expected token %d, got %d", typ, it.typ))
	}
}

type itemType int

const (
	itemEOF = iota
	itemNumber
	itemComma
	itemDash
)

type item struct {
	typ itemType
	s   string
}

type stateFn func(*scanner) stateFn

type scanner struct {
	input []byte
	start int
	pos   int
	item  item
	isEOF bool
}

func (s *scanner) nextItem() item {
	s.item = item{itemEOF, ""}
	state := scanInstruction
	for {
		state = state(s)
		if state == nil {
			return s.item
		}
	}
}

func (s *scanner) next() byte {
	if s.pos >= len(s.input) {
		s.isEOF = true
		return 0
	}
	n := s.input[s.pos]
	s.pos++
	return n
}

func (s *scanner) backup() {
	if s.isEOF {
		return
	}
	s.pos--
}

func (s *scanner) emit(typ itemType) stateFn {
	i := item{typ: typ, s: string(s.input[s.start:s.pos])}
	s.start = s.pos
	s.item = i
	return nil
}

func scanInstruction(s *scanner) stateFn {
	n := s.next()
	switch {
	case s.isEOF || n == '\n':
		return nil
	case n == '-':
		return s.emit(itemDash)
	case n == ',':
		return s.emit(itemComma)
	case isNumber(n):
		return scanNumber
	default:
		panic(fmt.Sprintf("unexpected char: %s", string(n)))
	}
}

func isNumber(b byte) bool { return b >= '0' && b <= '9' }

func scanNumber(s *scanner) stateFn {
	n := s.next()
	for isNumber(n) {
		n = s.next()
	}
	s.backup()
	return s.emit(itemNumber)
}
