package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	available, requested := parse(input)
	tr := newTrie()
	for _, d := range available {
		tr.add(d)
	}

	var count int
	for _, d := range requested {
		count += PathExists(tr, d)
	}
	fmt.Println(count)
}
