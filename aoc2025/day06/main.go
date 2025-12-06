package main

import (
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	star2(input)
}

func star1(input []byte) {
	g := parseGridStar1(string(input))
	var grandTotal int64

	for col := 0; col < g.width; col++ {
		op := g.op(col)
		var colTotal int64
		if op == "*" {
			colTotal = 1
		}
		for nStr := range g.iterNumbersCol(col) {
			n, err := strconv.ParseInt(nStr, 10, 64)
			if err != nil {
				panic(err)
			}

			switch op {
			case "*":
				colTotal *= n
			case "+":
				colTotal += n
			}
		}

		grandTotal += colTotal
	}

	fmt.Println(grandTotal)
}

type gridStar1 struct {
	rows   [][]string
	height int
	width  int
}

func (g gridStar1) get(x, y int) string {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return ""
	}
	return g.rows[y][x]
}

func (g gridStar1) op(col int) string {
	return g.get(col, g.height-1)
}

func (g gridStar1) iterNumbersCol(col int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for y := range g.height - 1 {
			if !yield(g.get(col, y)) {
				return
			}
		}
	}
}

func (g gridStar1) String() string {
	var sb strings.Builder
	for _, row := range g.rows {
		for _, cell := range row {
			sb.WriteString(cell)
			sb.WriteRune('\t')
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func parseGridStar1(in string) gridStar1 {
	lines := strings.Split(in, "\n")

	rows := make([][]string, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var row []string
		start := strings.IndexAny(line, "1234567890*+")
		for i := start; i < len(line); i++ {
			if line[i] == ' ' {
				row = append(row, line[start:i])
				for i < len(line) && line[i] == ' ' {
					i++
				}
				start = i
				continue
			}
		}
		if start < len(line) {
			row = append(row, line[start:])
		}
		rows = append(rows, row)
	}

	height := len(rows)
	var width int
	if height > 0 {
		width = len(rows[0])
	}

	// input check: all rows are same width
	for _, r := range rows {
		if len(r) != width {
			panic(fmt.Errorf("mismatched width: expected %d got %d", width, len(r)))
		}
	}

	return gridStar1{rows, height, width}
}

func star2(input []byte) {
	problems := parseStar2(string(input))

	var grandTotal int64
	for _, p := range problems {
		grandTotal += p.Result()
	}
	fmt.Println(grandTotal)
}

type Problem struct {
	op       rune
	nums     []int64
	colStart int
	colEnd   int
}

func (p Problem) Result() int64 {
	var res int64
	if p.op == '*' {
		res = 1
	}

	for _, n := range p.nums {
		switch p.op {
		case '+':
			res += n
		case '*':
			res *= n
		}
	}

	return res
}

func parseStar2(in string) []*Problem {
	lines := strings.Split(in, "\n")
	numLines := lines[:len(lines)-2]
	opLine := lines[len(lines)-2]

	var problems []*Problem
	for i, c := range opLine {
		if c != ' ' {
			if i != 0 {
				problems[len(problems)-1].colEnd--
			}
			problems = append(problems, &Problem{
				op:       c,
				colStart: i,
				colEnd:   i,
			})
		} else {
			problems[len(problems)-1].colEnd++
		}
	}

	for _, p := range problems {
		for col := p.colStart; col <= p.colEnd; col++ {
			var nStr string
			for row := range numLines {
				digit := rune(numLines[row][col])
				if digit == ' ' {
					continue
				}
				nStr += string(digit)
			}

			n, err := strconv.ParseInt(nStr, 10, 64)
			if err != nil {
				panic(err)
			}

			p.nums = append(p.nums, n)
		}
	}

	return problems
}
