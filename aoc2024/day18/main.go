package main

import (
	_ "embed"
	"fmt"
	"iter"
	"strconv"
	"strings"
)

type Grid [][]byte

func Build(size int, in []byte) (Grid, []vec) {
	g := make([][]byte, size)
	for i := range g {
		g[i] = make([]byte, size)
	}

	var (
		cur    int
		blocks []vec
	)
	for cur < len(in) {
		x := parseInt(&cur, in)
		if in[cur] != ',' {
			panic("expecting ','")
		}
		cur++

		y := parseInt(&cur, in)
		if cur < len(in) && in[cur] == '\n' {
			cur++
		}

		blocks = append(blocks, vec{x, y})
	}

	return g, blocks
}

func parseInt(cur *int, in []byte) int {
	start := *cur
	for *cur < len(in) && in[*cur] >= '0' && in[*cur] <= '9' {
		*cur++
	}
	v, err := strconv.Atoi(string(in[start:*cur]))
	if err != nil {
		panic(err)
	}
	return v
}

func (g Grid) Get(x, y int) byte {
	if x < 0 || y < 0 || y >= len(g) || x >= len(g[0]) {
		return '#'
	}
	return g[y][x]
}

func (g Grid) String() string {
	var b strings.Builder
	for _, row := range g {
		for _, ch := range row {
			if ch == 0 {
				b.WriteByte('.')
			} else {
				b.WriteByte(ch)
			}
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

func (g Grid) PrintPath(path []node) {
	set := make(set)
	for _, p := range path {
		set.add(p)
	}

	for y, row := range g {
		for x, ch := range row {
			if set.contains(node{x: x, y: y}) {
				fmt.Printf("O")
			} else if ch == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%c", ch)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type vec struct {
	x, y int
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

func isWall(b byte) bool {
	return b == '#'
}

type node struct {
	x, y int
}

func reconstructPath(cameFrom map[node]node, current node) []node {
	var totalPath []node
	totalPath = append(totalPath, current)
	for {
		c, found := cameFrom[current]
		if !found {
			break
		}
		current = c
		totalPath = append(totalPath, current)
	}
	return totalPath
}

type set map[node]struct{}

func (s set) add(n node) {
	s[n] = struct{}{}
}

func (s set) remove(n node) {
	delete(s, n)
}

func (s set) contains(n node) bool {
	_, ok := s[n]
	return ok
}

func (s set) len() int {
	return len(s)
}

type scoreMap map[vec]int

func (s scoreMap) get(x, y int) int {
	res, ok := s[vec{x, y}]
	if !ok {
		return 999999999999999999
	}
	return res
}

func (s scoreMap) set(x, y, score int) {
	s[vec{x, y}] = score
}

func astar2(g Grid, goalX, goalY int) []node {
	openSet := make(set)
	openSet.add(node{x: 0, y: 0})

	cameFrom := make(map[node]node)
	gScore := make(scoreMap)
	gScore.set(0, 0, 0)
	fScore := make(scoreMap)
	fScore.set(0, 0, h(0, 0, goalX, goalY))

	for openSet.len() > 0 {
		current := popMin(openSet, fScore)
		openSet.remove(current)
		if current.x == goalX && current.y == goalY {
			return reconstructPath(cameFrom, current)
		}

		neighbors := []node{
			{x: current.x + 1, y: current.y},
			{x: current.x, y: current.y + 1},
			{x: current.x - 1, y: current.y},
			{x: current.x, y: current.y - 1},
		}

		for _, n := range neighbors {
			if isWall(g.Get(n.x, n.y)) {
				continue
			}
			d := 1
			tentativeGScore := gScore.get(current.x, current.y) + d
			if tentativeGScore < gScore.get(n.x, n.y) {
				cameFrom[n] = current
				gScore.set(n.x, n.y, tentativeGScore)
				fScore.set(n.x, n.y, tentativeGScore+h(n.x, n.y, goalX, goalY))
				if !openSet.contains(n) {
					openSet.add(n)
				}
			}
		}
	}

	return nil
}

func h(ax, ay, bx, by int) int {
	absX := abs(ax - bx)
	absY := abs(ay - by)
	return absX + absY
}

func popMin(nodes set, scores scoreMap) node {
	if len(nodes) == 0 {
		panic("popMin: empty list")
	}

	var min node
	for v := range nodes {
		min = v
		break
	}

	for n := range nodes {
		if scores.get(n.x, n.y) < scores.get(min.x, min.y) {
			min = n
		}
	}

	return min
}

func heur(a node, b vec) int {
	startX := a.x
	startY := a.y
	endX := b.x
	endY := b.y

	if startX > endX {
		startX, endX = endX, startX
	}
	if startY > endY {
		startY, endY = endY, startY
	}

	m := startX - endX
	m += startY - endY

	return m
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//go:embed input.txt
var input []byte

func main() {
	// star1()
	star2()
}

func star1() {
	g, blocks := Build(71, input)
	blocks = blocks[:1024]
	for _, b := range blocks {
		g[b.y][b.x] = '#'
	}
	path := astar2(g, 70, 70)
	fmt.Println(len(path) - 1)
}

func star2() {
	_, blocks := Build(71, input)
	i := canBeSolvedBinarySearch(71, blocks, 0, len(blocks))
	fmt.Printf("%d,%d\n", blocks[i].x, blocks[i].y)
}

func canBeSolvedBinarySearch(size int, blocks []vec, start, end int) int {
	if end == start {
		return end - 1
	}

	g := make(Grid, size)
	for i := range size {
		g[i] = make([]byte, size)
	}

	middle := (start + end) / 2

	for _, b := range blocks[:middle] {
		g[b.y][b.x] = '#'
	}

	path := astar2(g, size-1, size-1)
	solvable := len(path) != 0

	if solvable {
		return canBeSolvedBinarySearch(size, blocks, middle+1, end)
	}
	return canBeSolvedBinarySearch(size, blocks, start, middle)
}
