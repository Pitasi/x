package main

import "fmt"

func RunStar1(g [][]*Node) int {
	values := findValues(g)
	var price int
	for _, val := range values {
		areas := findAreasForValue(g, val)
		for _, a := range areas {
			p := computePrice(a)
			price += p
		}
	}
	return price
}

func RunStar2(g [][]*Node) int {
	values := findValues(g)
	var price int
	for _, val := range values {
		areas := findAreasForValue(g, val)
		for _, a := range areas {
			p := computePrice2(a)
			price += p
		}
	}
	return price
}

type Node struct {
	Value byte
	Up    *Node
	Down  *Node
	Left  *Node
	Right *Node

	X, Y int
}

func (n *Node) GetValue() byte {
	if n == nil {
		return 0
	}
	return n.Value
}

func (n *Node) String() string {
	return fmt.Sprintf("{%c,%p}", n.Value, n)
}

func findValues(grid [][]*Node) []byte {
	result := make(map[byte]struct{})
	for _, row := range grid {
		for _, node := range row {
			if _, ok := result[node.Value]; !ok {
				result[node.Value] = struct{}{}
			}
		}
	}
	resultSlice := make([]byte, 0, len(result))
	for k := range result {
		resultSlice = append(resultSlice, k)
	}
	return resultSlice
}

func findAreasForValue(grid [][]*Node, value byte) [][]*Node {
	visited := make(set)
	var areas [][]*Node

	for _, row := range grid {
		for _, node := range row {
			if node.Value == value && !visited.has(node) {
				areas = append(areas, findContigArea(node, visited))
			}
		}
	}

	return areas
}

type set map[*Node]struct{}

func (s set) add(n *Node) {
	s[n] = struct{}{}
}

func (s set) has(n *Node) bool {
	_, ok := s[n]
	return ok
}

func findContigArea(n *Node, visited set) []*Node {
	area := []*Node{}
	if !visited.has(n) {
		area = append(area, n)
		visited.add(n)
	}

	neighbors := findNeighbors(n, n.Value, visited)
	for _, child := range neighbors {
		childNodes := findContigArea(child, visited)
		area = append(area, childNodes...)
	}

	return area
}

func findNeighbors(n *Node, searchValue byte, visited set) []*Node {
	if n == nil {
		return nil
	}
	result := make([]*Node, 0, 4)
	for _, child := range []*Node{n.Up, n.Down, n.Left, n.Right} {
		if child == nil || visited.has(child) {
			continue
		}
		if child.Value == searchValue {
			result = append(result, child)
		}
	}
	return result
}

func computePrice(area []*Node) int {
	a := computeAreaSize(area)
	p := computeAreaPerimeter(area)
	return a * p
}

func computePrice2(area []*Node) int {
	a := computeAreaSize(area)
	p := computeAreaSides(area)
	return a * p
}

func computeAreaSize(area []*Node) int {
	return len(area)
}

func computeAreaPerimeter(area []*Node) int {
	var p int
	for _, n := range area {
		for _, c := range []*Node{n.Up, n.Down, n.Left, n.Right} {
			if c == nil || c.Value != n.Value {
				p += 1
			}
		}
	}
	return p
}

func computeAreaSides(area []*Node) int {

	const (
		up = iota
		down
		left
		right
	)
	type k struct {
		node *Node
		dir  int
	}
	visited := make(map[k]struct{})
	has := func(n *Node, dir int) bool {
		_, ok := visited[k{n, dir}]
		return ok
	}

	var count int

	for _, n := range area {
		if n.Up != nil && n.Down != nil && n.Left != nil && n.Right != nil && n.Up.Value == n.Value && n.Down.Value == n.Value && n.Left.Value == n.Value && n.Right.Value == n.Value {
			// node in the middle of the area
			visited[k{n, up}] = struct{}{}
			visited[k{n, down}] = struct{}{}
			visited[k{n, left}] = struct{}{}
			visited[k{n, right}] = struct{}{}
			continue
		}

		if !has(n, up) && (n.Up == nil || n.Up.Value != n.Value) {
			curr := n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Up.GetValue() != n.GetValue() {
				visited[k{curr, up}] = struct{}{}
				curr = curr.Left
			}
			curr = n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Up.GetValue() != n.GetValue() {
				visited[k{curr, up}] = struct{}{}
				curr = curr.Right
			}
			count++
		}

		if !has(n, down) && (n.Down == nil || n.Down.Value != n.Value) {
			curr := n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Down.GetValue() != n.GetValue() {
				visited[k{curr, down}] = struct{}{}
				curr = curr.Left
			}
			curr = n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Down.GetValue() != n.GetValue() {
				visited[k{curr, down}] = struct{}{}
				curr = curr.Right
			}
			count++
		}

		if !has(n, left) && (n.Left == nil || n.Left.Value != n.Value) {
			curr := n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Left.GetValue() != n.GetValue() {
				visited[k{curr, left}] = struct{}{}
				curr = curr.Up
			}
			curr = n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Left.GetValue() != n.GetValue() {
				visited[k{curr, left}] = struct{}{}
				curr = curr.Down
			}
			count++
		}

		if !has(n, right) && (n.Right == nil || n.Right.Value != n.Value) {
			curr := n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Right.GetValue() != n.GetValue() {
				visited[k{curr, right}] = struct{}{}
				curr = curr.Up
			}
			curr = n
			for curr != nil && curr.GetValue() == n.GetValue() && curr.Right.GetValue() != n.GetValue() {
				visited[k{curr, right}] = struct{}{}
				curr = curr.Down
			}
			count++
		}
	}

	return count
}
