package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	m, robot, moves := Parse(input)
	score := Run(m, robot, moves)
	fmt.Println(score)

	m, robot, moves = Parse(input)
	mapv2, robot := MapV2(m)
	score2 := RunV2(mapv2, robot, moves)
	fmt.Println(score2)
}
