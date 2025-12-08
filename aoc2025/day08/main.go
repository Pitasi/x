package main

import (
	"errors"
	"fmt"
	"math"
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
	c := parseInput(input)

	var boxes []*box
	for _, a := range c {
		boxes = append(boxes, &box{coords: a})
	}

	var pairs []pair
	for i, a := range boxes {
		for _, b := range boxes[i+1:] {
			p := pair{
				a:    a,
				b:    b,
				dist: distance(a.coords, b.coords),
			}
			pairs = append(pairs, p)
		}
	}

	slices.SortFunc(pairs, func(a, b pair) int {
		return int(a.dist - b.dist)
	})

	pairs = pairs[:1000]
	for _, p := range pairs {
		connect(p.a, p.b)
	}

	var circuits []int
	for _, b := range boxes {
		if b.visited {
			continue
		}
		circuit := make(map[*box]struct{})
		visit(b, circuit)

		for b := range circuit {
			b.visited = true
		}

		circuits = append(circuits, len(circuit))
	}

	slices.Sort(circuits)
	slices.Reverse(circuits)

	res := circuits[0] * circuits[1] * circuits[2]

	fmt.Println(res)
}

func star2(input []byte) {
	c := parseInput(input)

	var boxes []*box
	for _, a := range c {
		boxes = append(boxes, &box{coords: a})
	}

	var pairs []pair
	for i, a := range boxes {
		for _, b := range boxes[i+1:] {
			p := pair{
				a:    a,
				b:    b,
				dist: distance(a.coords, b.coords),
			}
			pairs = append(pairs, p)
		}
	}

	slices.SortFunc(pairs, func(a, b pair) int {
		return int(a.dist - b.dist)
	})

	for _, p := range pairs {
		connect(p.a, p.b)
		circuit := make(map[*box]struct{})
		visit(boxes[0], circuit)
		if len(circuit) == len(boxes) {
			fmt.Println(p.a.coords.x * p.b.coords.x)
			break
		}
	}
}

func visit(b *box, visited map[*box]struct{}) {
	if _, found := visited[b]; found {
		return
	}

	visited[b] = struct{}{}
	for _, neigh := range b.conns {
		visit(neigh, visited)
	}
}

type pair struct {
	a    *box
	b    *box
	dist float64
}

func distance(a, b coords) float64 {
	return math.Sqrt(
		math.Pow(float64(a.x-b.x), 2) +
			math.Pow(float64(a.y-b.y), 2) +
			math.Pow(float64(a.z-b.z), 2),
	)
}

type box struct {
	coords coords
	conns  []*box

	visited bool
}

func connect(a, b *box) {
	a.conns = append(a.conns, b)
	b.conns = append(b.conns, a)
}

type coords struct {
	x int
	y int
	z int
}

func parseInput(input []byte) []coords {
	s := &scanner{input: input}

	var c []coords
	for {
		xStr := s.nextItem()
		if xStr.typ == itemEOF {
			break
		}
		if xStr.typ != itemNumber {
			panic("expected number")
		}
		comma1 := s.nextItem()
		if comma1.typ != itemComma {
			panic("expected comma")
		}
		yStr := s.nextItem()
		if yStr.typ != itemNumber {
			panic("expected number")
		}
		comma2 := s.nextItem()
		if comma2.typ != itemComma {
			panic("expected comma")
		}
		zStr := s.nextItem()
		if zStr.typ != itemNumber {
			panic("expected number")
		}
		nlOrEOF := s.nextItem()
		if nlOrEOF.typ == itemEOF {
			break
		}
		if nlOrEOF.typ != itemNewline {
			panic("expected comma or newline")
		}

		x, y, z, err := parseInts(xStr.s, yStr.s, zStr.s)
		if err != nil {
			panic(err)
		}

		c = append(c, coords{x, y, z})
	}

	return c
}

func parseInts(xS, yS, zS string) (int, int, int, error) {
	x, xErr := strconv.Atoi(xS)
	y, yErr := strconv.Atoi(yS)
	z, zErr := strconv.Atoi(zS)
	return x, y, z, errors.Join(xErr, yErr, zErr)
}

type itemType int

const (
	itemEOF = iota
	itemNumber
	itemComma
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
	case n == ',':
		return s.emit(itemComma)
	case isNumber(n):
		return scanNumber
	}
	panic("unexpected token")
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
