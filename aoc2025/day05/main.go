package main

import (
	"fmt"
	"os"
	"slices"
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

func star1(input []byte) {
	s := &scanner{input: input}

	var ranges []Range
	for {
		r, ok := parseRange(s)
		if !ok {
			break
		}
		ranges = append(ranges, r)
	}

	slices.SortFunc(ranges, func(a, b Range) int {
		return a.start - b.start
	})

	var count int
	for {
		n, ok := parseNumber(s)
		if !ok {
			break
		}

		for _, r := range ranges {
			if r.contains(n) {
				count++
				break
			}
		}
	}

	fmt.Println(count)
}

func star2(input []byte) {
	s := &scanner{input: input}

	var ranges []Range
	for {
		r, ok := parseRange(s)
		if !ok {
			break
		}
		ranges = append(ranges, r)
	}

	slices.SortFunc(ranges, func(a, b Range) int {
		return a.start - b.start
	})

	for i := 0; i < len(ranges)-1; i++ {
		curr := ranges[i]
		if curr.empty {
			continue
		}
		for j := i + 1; j < len(ranges); j++ {
			if curr.end >= ranges[j].start {
				ranges[j].start = curr.end + 1
				if ranges[j].start > ranges[j].end {
					ranges[j].empty = true
				} else {
					break
				}
			}
		}
	}

	var count int
	for _, r := range ranges {
		count += r.len()
	}
	fmt.Println(count)
}

type Range struct {
	start int
	end   int
	empty bool
}

func (r Range) contains(n int) bool { return r.start <= n && r.end >= n }

func (r Range) len() int {
	if r.empty {
		return 0
	}
	return r.end - r.start + 1
}

func parseRange(s *scanner) (Range, bool) {
	rangeStart := s.nextItem()
	if rangeStart.typ == itemNL {
		return Range{}, false
	}

	if rangeStart.typ != itemNumber {
		panic("expected number")
	}

	dash := s.nextItem()
	if dash.typ != itemDash {
		panic("expected dash")
	}

	rangeEnd := s.nextItem()
	if rangeEnd.typ != itemNumber {
		panic("expected number")
	}

	nl := s.nextItem()
	if nl.typ != itemNL {
		panic("expected new line")
	}

	rs, _ := strconv.ParseInt(rangeStart.s, 10, 64)
	re, _ := strconv.ParseInt(rangeEnd.s, 10, 64)
	return Range{start: int(rs), end: int(re)}, true
}

func parseNumber(s *scanner) (int, bool) {
	n := s.nextItem()
	if n.typ == itemEOF || n.typ == itemNL {
		return 0, false
	}
	if n.typ != itemNumber {
		panic("expected number")
	}

	nl := s.nextItem()
	if nl.typ != itemEOF && nl.typ != itemNL {
		panic("expected new line or EOF")
	}

	nn, _ := strconv.Atoi(n.s)
	return nn, true
}

type itemType int

const (
	itemEOF = iota
	itemNumber
	itemDash
	itemNL
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
		return s.emit(itemNL)
	case n == '-':
		return s.emit(itemDash)
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
