package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	grid := Parse(input)
	roots := Roots(grid)

	var totalScore int
	for _, root := range roots {
		trails := HikingTrails(root, 9)
		score := CountUniques(trails)
		totalScore += score
	}
	fmt.Println(totalScore)

	var totalRating int
	for _, root := range roots {
		trails := HikingTrails(root, 9)
		totalRating += len(trails)
	}
	fmt.Println(totalRating)
}
