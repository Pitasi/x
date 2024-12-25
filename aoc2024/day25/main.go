package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	keys := Parse(input)
	fmt.Println(star1(keys))
}

func star1(keys []Key) int {
	var count int
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i].Fit(keys[j]) {
				count++
			}
		}
	}
	return count
}
