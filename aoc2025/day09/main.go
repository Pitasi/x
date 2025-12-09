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

	s := &scanner{input: input}

	var coords []Coords
	for {
		r, ok := parseCoords(s)
		if !ok {
			break
		}
		coords = append(coords, r)
	}

	star1(coords)
	star2(coords)
}

func star1(coords []Coords) {
	var m int
	for i, a := range coords {
		for _, b := range coords[i+1:] {
			ar := area(a, b)
			m = max(ar, m)
		}
	}

	fmt.Println(m)
}

func star2(coords []Coords) {
	var m int

	coords = append(coords, coords[len(coords)-1]) // close polygon

	for i, a := range coords {
		for _, b := range coords[i+1:] {
			minX := min(a.x, b.x)
			minY := min(a.y, b.y)
			maxX := max(a.x, b.x)
			maxY := max(a.y, b.y)

			if !intersects(minX, minY, maxX, maxY, coords) {
				ar := area(a, b)
				if ar > m {
					m = ar
				}
			}
		}
	}

	fmt.Println(m)
}

func intersects(minX, minY, maxX, maxY int, coords []Coords) bool {
	for i := 0; i < len(coords)-1; i++ {
		a := coords[i]
		b := coords[i+1]

		boundLeft := min(a.x, b.x)
		boundRight := max(a.x, b.x)
		boundTop := min(a.y, b.y)
		boundBottom := max(a.y, b.y)

		if minX < boundRight && maxX > boundLeft && minY < boundTop && maxY > boundBottom {
			return true
		}
	}

	return false
}

func area(a, b Coords) int {
	h := abs(a.y-b.y) + 1
	if a.x == b.x {
		return h
	}

	w := abs(a.x-b.x) + 1
	if a.y == b.y {
		return w
	}

	return h * w
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type Coords struct {
	x int
	y int
}

func parseCoords(s *scanner) (Coords, bool) {
	rangeStart := s.nextItem()
	if rangeStart.typ == itemNL || rangeStart.typ == itemEOF {
		return Coords{}, false
	}
	if rangeStart.typ != itemNumber {
		panic("expected number")
	}

	dash := s.nextItem()
	if dash.typ != itemComma {
		panic("expected comma")
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
	return Coords{x: int(rs), y: int(re)}, true
}

type itemType int

const (
	itemEOF = iota
	itemNumber
	itemComma
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
