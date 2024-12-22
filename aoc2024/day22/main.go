package main

import (
	_ "embed"
	"fmt"
)

func random(seed int) int {
	seed ^= seed << 6
	seed &= 0b111111111111111111111111
	seed ^= seed >> 5
	seed &= 0b111111111111111111111111
	seed ^= seed << 11
	seed &= 0b111111111111111111111111
	return seed
}

func mix(v, seed int) int {
	return seed ^ v
}

func nth(seed int, n int) int {
	for i := 0; i < n; i++ {
		seed = random(seed)
	}
	return seed
}

//go:embed input.txt
var input []byte

type vec struct {
	a, b, c, d int
}

type buyer struct {
	scores map[vec]int
}

func newBuyer(seed int) *buyer {
	seq := make([]int, 0, 2000)
	seq = append(seq, seed)
	for i := 1; i < 2000; i++ {
		seq = append(seq, random(seq[i-1]))
	}

	for i := 0; i < len(seq); i++ {
		seq[i] %= 10
	}

	diffs := make([]int, len(seq))
	diffs[0] = seq[0] % 10
	for i := 1; i < len(seq); i++ {
		diffs[i] = seq[i] - seq[i-1]
	}

	scores := make(map[vec]int)
	for i := 4; i < len(diffs); i++ {
		a := diffs[i-3]
		b := diffs[i-2]
		c := diffs[i-1]
		d := diffs[i]
		key := vec{a, b, c, d}
		if _, ok := scores[key]; ok {
			continue
		}
		scores[key] = seq[i]
	}

	return &buyer{scores: scores}
}

func allPossibleSeqs(buyers []*buyer) int {
	set := make(map[vec]int)
	for _, buyer := range buyers {
		for key, v := range buyer.scores {
			set[key] += v
		}
	}

	var m int
	for _, v := range set {
		m = max(m, v)
	}

	return m
}

func main() {
	seeds := Parse(input)
	buyers := make([]*buyer, len(seeds))
	for i, seed := range seeds {
		buyers[i] = newBuyer(seed)
	}

	fmt.Println(allPossibleSeqs(buyers))
}

func star1() {
	seeds := Parse(input)

	var sum int
	for _, seed := range seeds {
		sum += nth(seed, 2000)
	}

	fmt.Println(sum)
}
