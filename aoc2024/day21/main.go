package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"strconv"
)

//go:embed input.txt
var input []byte

func main() {
	star := day21(input, 24)
	fmt.Println(star)
}

func day21(input []byte, robots int) int {
	lines := bytes.Split(input, []byte("\n"))
	lookup := make(map[string]int)

	var sum int
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		c := digits(line)
		s := npad(line, robots, lookup)
		sum += c * s
	}
	return sum
}

func digits(input []byte) int {
	v, err := strconv.Atoi(string(input[:3]))
	if err != nil {
		panic(err)
	}
	return v
}

func npad(input []byte, robots int, lookup map[string]int) int {
	input = append([]byte("A"), input...)
	var acc int
	for i := 1; i < len(input); i++ {
		start := input[i-1]
		end := input[i]
		sequences := numpad(start, end)
		results := make([]int, 0, len(sequences))
		for _, seq := range sequences {
			results = append(results, dpad(seq, robots, lookup))
		}
		acc += slices.Min(results)
	}

	return acc
}

func dpad(input []byte, robots int, lookup map[string]int) int {
	key := strconv.Itoa(robots) + string(input)
	if cached, found := lookup[key]; found {
		return cached
	}

	input = append([]byte("A"), input...)
	var acc int
	for i := 1; i < len(input); i++ {
		start := input[i-1]
		end := input[i]
		sequences := arrowpad(start, end)

		results := make([]int, 0, len(sequences))
		if robots == 0 {
			for _, seq := range sequences {
				results = append(results, len(seq))
			}
		} else {
			for _, seq := range sequences {
				results = append(results, dpad(seq, robots-1, lookup))
			}
		}
		acc += slices.Min(results)
	}

	lookup[key] = acc
	return acc
}

func numpad(start, end byte) [][]byte {
	startX, startY := numpadCoords(start)
	endX, endY := numpadCoords(end)

	dX := endX - startX
	xByte := byte('>')
	if dX < 0 {
		xByte = '<'
		dX = -dX
	}
	xBytes := make([]byte, dX)
	for i := 0; i < dX; i++ {
		xBytes[i] = xByte
	}

	dY := endY - startY
	yByte := byte('v')
	if dY < 0 {
		yByte = '^'
		dY = -dY
	}
	yBytes := make([]byte, dY)
	for i := 0; i < dY; i++ {
		yBytes[i] = yByte
	}

	var res [][]byte

	if !(startY == 3 && endX == 0) {
		combination := append(xBytes, yBytes...)
		combination = append(combination, 'A')
		res = append(res, combination)
	}

	if !(startX == 0 && endY == 3) {
		combination := append(yBytes, xBytes...)
		combination = append(combination, 'A')
		if len(res) == 0 || !bytes.Equal(combination, res[0]) {
			res = append(res, combination)
		}
	}

	return res
}

func arrowpad(start, end byte) [][]byte {
	startX, startY := arrowpadCoords(start)
	endX, endY := arrowpadCoords(end)

	dX := endX - startX
	xByte := byte('>')
	if dX < 0 {
		xByte = '<'
		dX = -dX
	}
	xBytes := make([]byte, dX)
	for i := 0; i < dX; i++ {
		xBytes[i] = xByte
	}

	dY := endY - startY
	yByte := byte('v')
	if dY < 0 {
		yByte = '^'
		dY = -dY
	}
	yBytes := make([]byte, dY)
	for i := 0; i < dY; i++ {
		yBytes[i] = yByte
	}

	var res [][]byte

	if !(startY == 0 && endX == 0) {
		combination := append(xBytes, yBytes...)
		combination = append(combination, 'A')
		res = append(res, combination)
	}

	if !(startX == 0 && endY == 0) {
		combination := append(yBytes, xBytes...)
		combination = append(combination, 'A')
		if len(res) == 0 || !bytes.Equal(combination, res[0]) {
			res = append(res, combination)
		}
	}

	return res
}

func arrowpadCoords(key byte) (int, int) {
	switch key {
	case '^':
		return 1, 0
	case 'A':
		return 2, 0
	case '<':
		return 0, 1
	case 'v':
		return 1, 1
	case '>':
		return 2, 1
	}

	panic("invalid key")
}

func numpadCoords(key byte) (int, int) {
	switch key {
	case '7':
		return 0, 0
	case '8':
		return 1, 0
	case '9':
		return 2, 0
	case '4':
		return 0, 1
	case '5':
		return 1, 1
	case '6':
		return 2, 1
	case '1':
		return 0, 2
	case '2':
		return 1, 2
	case '3':
		return 2, 2
	case '0':
		return 1, 3
	case 'A':
		return 2, 3
	}

	panic("invalid key")
}
