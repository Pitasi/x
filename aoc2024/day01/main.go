package main

import (
	_ "embed"
	"fmt"
	"sort"
)

func main() {
	t := Tokenizer{input: input}
	l1 := make([]int, 0, 10000)
	l2 := make([]int, 0, 10000)

	for {
		l, r, ok := t.Next()
		if !ok {
			break
		}
		l1 = append(l1, l)
		l2 = append(l2, r)
	}
	if len(l1) != len(l2) {
		panic("lists have different lengths")
	}

	sort.Ints(l1)
	sort.Ints(l2)

	firstStar(l1, l2)
	secondStar(l1, l2)
}

func firstStar(l1, l2 []int) {
	var diff int
	for i := 0; i < len(l1); i++ {
		d := l1[i] - l2[i]
		if d < 0 {
			d = -d
		}
		diff += d
	}

	fmt.Println(diff)
}

func secondStar(l1, l2 []int) {
	scores := make(map[int]int)

	last := l2[0]
	lastScore := 1
	for i := 1; i < len(l2); i++ {
		if l2[i] == last {
			lastScore++
		} else {
			scores[last] = lastScore
			last = l2[i]
			lastScore = 1
		}
	}

	var totalScore int
	for _, n := range l1 {
		totalScore += n * scores[n]
	}

	fmt.Println(totalScore)
}
