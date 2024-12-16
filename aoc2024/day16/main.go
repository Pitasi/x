package main

import (
	_ "embed"
	"fmt"
	"iter"
	"strings"
)

type Dir int

const (
	North Dir = iota
	East
	South
	West
)

func Parse(input []byte) (Grid, int, int, int, int) {
	g := make(Grid, 0)
	var startX, startY, goalX, goalY int
	var currentRow []byte

	for i, ch := range input {
		x := len(currentRow)
		y := len(g)

		switch input[i] {
		case 'S':
			startX = x
			startY = y
			currentRow = append(currentRow, '.')
		case 'E':
			goalX = x
			goalY = y
			currentRow = append(currentRow, '.')
		case '\n':
			g = append(g, currentRow)
			currentRow = nil
		default:
			currentRow = append(currentRow, ch)
		}
	}
	if len(currentRow) > 0 {
		g = append(g, currentRow)
	}

	return g, startX, startY, goalX, goalY
}

type Grid [][]byte

func (g Grid) Get(x, y int) byte {
	if x < 0 || y < 0 || y >= len(g) || x >= len(g[0]) {
		return 0
	}
	return g[y][x]
}

func (g Grid) String() string {
	var b strings.Builder
	for _, row := range g {
		for _, ch := range row {
			b.WriteByte(ch)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (g Grid) Print(sX, sY, eX, eY int) {
	for y, row := range g {
		for x, ch := range row {
			if x == sX && y == sY {
				fmt.Printf(" S")
			} else if x == eX && y == eY {
				fmt.Printf(" E")
			} else {
				fmt.Printf(" %c", ch)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g Grid) PrintPath(p path) {
	for y, row := range g {
		for x, ch := range row {
			if p.Contains(x, y) {
				fmt.Printf(" O")
			} else {
				fmt.Printf(" %c", ch)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type vec struct {
	x, y int
}

type scores map[vec]int

func (s scores) Get(x, y int) int {
	res, ok := s[vec{x, y}]
	if !ok {
		return 999999999999999999
	}
	return res
}

func (s scores) Set(x, y, score int) {
	if score > s.Get(x, y) {
		return
	}
	s[vec{x, y}] = score
}

type set map[vec]struct{}

func (s set) Has(x, y int) bool {
	_, ok := s[vec{x, y}]
	return ok
}

func (s set) Add(x, y int) {
	s[vec{x, y}] = struct{}{}
}

type path []vec

func (p path) Add(x, y int) path {
	return append(p, vec{x, y})
}

func (p path) Last() (int, int) {
	if len(p) == 0 {
		return 0, 0
	}
	return p[len(p)-1].x, p[len(p)-1].y
}

func (p path) Contains(x, y int) bool {
	for _, it := range p {
		if it.x == x && it.y == y {
			return true
		}
	}
	return false
}

func (p path) Iter() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for _, it := range p {
			if !yield(it.x, it.y) {
				return
			}
		}
	}
}

func Run(g Grid, startX, startY, goalX, goalY int) (int, int) {
	sc := make(scores)
	finalScores := make(scores)
	score := 0
	dir := East
	maxSteps := countNonWalls(g)
	var currentPath path
	score, _ = f(g, startX, startY, currentPath, dir, goalX, goalY, score, -1, maxSteps, sc, finalScores)

	var count int
	for _, sc := range finalScores {
		if sc == score {
			count++
		}
	}

	return score, count
}

func countNonWalls(g Grid) int {
	count := 0
	for _, row := range g {
		for _, ch := range row {
			if !isWall(ch) {
				count++
			}
		}
	}
	return count
}

func turns(start, end Dir) int {
	if start == end {
		return 0
	}
	if (start == North && end == South) ||
		(end == North && start == South) ||
		(start == East && end == West) ||
		(start == West && end == East) {
		return 2
	}
	return 1
}

func isWall(b byte) bool {
	return b == 0 || b == '#'
}

func f(g Grid, x, y int, currentPath path, dir Dir, goalX, goalY int, score, currMinScore int, maxSteps int, scores, finalScores scores) (int, bool) {
	if len(currentPath) > maxSteps || (currMinScore > -1 && score > currMinScore) {
		return -1, false
	}

	tolerance := 1000
	if scores.Get(x, y) < score-tolerance {
		return -1, false
	}
	if scores.Get(x, y) > score {
		scores.Set(x, y, score)
	}
	currentPath = currentPath.Add(x, y)

	if x == goalX && y == goalY {
		for x, y := range currentPath.Iter() {
			finalScores.Set(x, y, score)
		}
		return score, true
	}

	up := g.Get(x, y-1)
	right := g.Get(x+1, y)
	down := g.Get(x, y+1)
	left := g.Get(x-1, y)

	minScore := -1

	if !isWall(up) {
		upScore := 1000*turns(dir, North) + 1 + score
		upScore, finished := f(g, x, y-1, currentPath, North, goalX, goalY, upScore, currMinScore, maxSteps, scores, finalScores)
		if finished && (currMinScore == -1 || upScore <= currMinScore) {
			currMinScore = upScore
		}
		if upScore != -1 && (minScore == -1 || upScore < minScore) {
			minScore = upScore
		}
	}

	if !isWall(right) {
		rightScore := 1000*turns(dir, East) + 1 + score
		rightScore, finished := f(g, x+1, y, currentPath, East, goalX, goalY, rightScore, currMinScore, maxSteps, scores, finalScores)
		if finished && (currMinScore == -1 || rightScore <= currMinScore) {
			currMinScore = rightScore
		}
		if rightScore != -1 && (minScore == -1 || rightScore < minScore) {
			minScore = rightScore
		}
	}

	if !isWall(down) {
		downScore := 1000*turns(dir, South) + 1 + score
		downScore, finished := f(g, x, y+1, currentPath, South, goalX, goalY, downScore, currMinScore, maxSteps, scores, finalScores)
		if finished && (currMinScore == -1 || downScore <= currMinScore) {
			currMinScore = downScore
		}
		if downScore != -1 && (minScore == -1 || downScore < minScore) {
			minScore = downScore
		}
	}

	if !isWall(left) {
		leftScore := 1000*turns(dir, West) + 1 + score
		leftScore, finished := f(g, x-1, y, currentPath, West, goalX, goalY, leftScore, currMinScore, maxSteps, scores, finalScores)
		if finished && (currMinScore == -1 || leftScore <= currMinScore) {
			currMinScore = leftScore
		}
		if leftScore != -1 && (minScore == -1 || leftScore < minScore) {
			minScore = leftScore
		}
	}

	return minScore, false
}

//go:embed input.txt
var input []byte

func main() {
	g, startX, startY, goalX, goalY := Parse(input)
	score, count := Run(g, startX, startY, goalX, goalY)
	fmt.Println(score, count)
}
