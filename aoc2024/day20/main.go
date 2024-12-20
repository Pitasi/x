package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	g := Parse(input)
	cands := CheatV2Candidates(g, 100, 20)
	fmt.Println(len(cands))
}

func star1() {
	g := Parse(input)
	candsH, candsV := CheatV1Candidates(g)
	cutoff := 99
	var count int
	for _, c := range candsH {
		s := CheatV1SavingH(g, c)
		if s > cutoff {
			count++
		}
	}
	for _, c := range candsV {
		s := CheatV1SavingV(g, c)
		if s > cutoff {
			count++
		}
	}
	fmt.Println(count)

	// or, using V2:
	// g := Parse(input)
	// cands := CheatV2Candidates(g, 100, 20)
	// fmt.Println(len(cands))
}
