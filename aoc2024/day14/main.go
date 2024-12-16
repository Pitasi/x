package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	robots := Parse(input)
	s1 := Star1(robots)
	fmt.Println(s1)

	m := v(101, 103)
	var count int
	for {
		count++
		robots = StepMultiple(1, robots, m)
		printed := Print(robots, m)
		if printed {
			fmt.Println("Seconds:", count)
			fmt.Println("Press ctrl+c to exit, or Enter to continue")
			fmt.Println("------------------------------------------------------------------")
			fmt.Scanln()
		}
	}
}

func Print(robots []Robot, mapSize Vector) bool {
	s := make([][]rune, mapSize.Y)
	for y := range s {
		s[y] = make([]rune, mapSize.X)
	}

	for _, robot := range robots {
		s[robot.P.Y][robot.P.X] = 'X'
	}

	if !containsRect(s, 3, 3) {
		return false
	}

	for y := range mapSize.Y {
		for x := range mapSize.X {
			if s[y][x] == 'X' {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
	return true
}

func containsRect(grid [][]rune, w, h int) bool {
	for y := range grid {
		for x := range grid[y] {
			if isRect(grid, x, y, w, h) {
				fmt.Println("yes", x, y, w, h)
				return true
			}
		}
	}
	return false
}

func isRect(grid [][]rune, x, y, w, h int) bool {
	x2 := x + w
	y2 := y + h
	for b := y; b < y2; b++ {
		for a := x; a < x2; a++ {
			if get(grid, a, b) != 'X' {
				return false
			}
		}

	}
	return true
}

func get(grid [][]rune, x, y int) rune {
	if y < 0 || y >= len(grid) {
		return 0
	}
	if x < 0 || x >= len(grid[0]) {
		return 0
	}
	return grid[y][x]
}
