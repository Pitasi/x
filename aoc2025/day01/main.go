package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	star1(input)
	star2(input)
}

type dial struct {
	n int64
}

func (d *dial) rotate(direction string, amount string) {
	amt, _ := strconv.ParseInt(amount, 10, 64)
	switch direction {
	case "L":
		d.n -= amt
	case "R":
		d.n += amt
	}
	d.n %= 100
}

func star1(input []byte) {
	s := &scanner{input: input}
	dial := &dial{n: 50}
	var count int

	for {
		dir := s.nextItem()
		if dir.typ == itemEOF {
			break
		} else if dir.typ != itemDirection {
			panic(fmt.Sprintf("expected direction got: %v", dir))
		}

		amt := s.nextItem()
		if amt.typ != itemNumber {
			panic(fmt.Sprintf("expected number got: %v", amt))
		}

		dial.rotate(dir.s, amt.s)
		if dial.n == 0 {
			count++
		}

		delim := s.nextItem()
		if delim.typ == itemEOF {
			break
		} else if delim.typ != itemNewline {
			panic(fmt.Sprintf("expected new line got: %v", delim))
		}
	}

	fmt.Println(count)
}

type dial2 struct {
	n     int64
	count int64
}

func (d *dial2) rotate(direction string, amount string) {
	amt, _ := strconv.ParseInt(amount, 10, 64)

	old := d.n
	switch direction {
	case "L":
		d.n -= amt
		if old != 0 && d.n <= 0 {
			d.count++
		}
		if d.n < -99 {
			d.count += -d.n / 100
		}
	case "R":
		d.n += amt
		if d.n >= 100 {
			d.count += d.n / 100
		}
	}

	d.n = mod(d.n, 100)
}

func mod(a, b int64) int64 {
	return (a%b + b) % b
}

func star2(input []byte) {
	s := &scanner{input: input}
	dial := &dial2{n: 50}

	for {
		dir := s.nextItem()
		if dir.typ == itemEOF {
			break
		} else if dir.typ != itemDirection {
			panic(fmt.Sprintf("expected direction got: %v", dir))
		}

		amt := s.nextItem()
		if amt.typ != itemNumber {
			panic(fmt.Sprintf("expected number got: %v", amt))
		}

		dial.rotate(dir.s, amt.s)

		delim := s.nextItem()
		if delim.typ == itemEOF {
			break
		} else if delim.typ != itemNewline {
			panic(fmt.Sprintf("expected new line got: %v", delim))
		}
	}

	fmt.Println(dial.count)
}

type Direction int

const (
	L Direction = iota
	R
)

type itemType int

const (
	itemEOF = iota
	itemDirection
	itemNumber
	itemNewline
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
	case s.isEOF:
		return nil
	case n == '\n':
		return s.emit(itemNewline)
	case isNumber(n):
		return scanNumber
	default:
		s.backup()
		return scanDirection
	}
}

func scanNumber(s *scanner) stateFn {
	n := s.next()
	for isNumber(n) {
		n = s.next()
	}
	s.backup()
	return s.emit(itemNumber)
}

func isNumber(b byte) bool { return b >= '0' && b <= '9' }

func scanDirection(s *scanner) stateFn {
	n := s.next()
	if n != 'L' && n != 'R' {
		panic("invalid direction: " + string(n))
	}
	return s.emit(itemDirection)
}
