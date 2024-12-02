package main

import "fmt"

func main() {
	t := Tokenizer{input: input}
	var safeCount int
	for {
		ints, ok := t.Next()
		if !ok {
			break
		}
		if isSafe(ints) || isSafeProblemDampener(ints) {
			safeCount++
		}
	}

	fmt.Println(safeCount)
}

func isSafeProblemDampener(ints []int) bool {
	for i := range ints {
		newReport := make([]int, len(ints)-1)
		// remove i from newReport
		copy(newReport, ints[:i])
		copy(newReport[i:], ints[i+1:])

		if isSafe(newReport) {
			return true
		}
	}
	return false
}

func isSafe(ints []int) bool {
	if !isSortedInc(ints) && !isSortedDec(ints) {
		return false
	}

	for i := 1; i < len(ints); i++ {
		diff := ints[i-1] - ints[i]
		if diff < 0 {
			diff = -diff
		}
		if diff < 1 || diff > 3 {
			return false
		}
	}

	return true
}

func isSortedInc(ints []int) bool {
	for i := 1; i < len(ints); i++ {
		if ints[i-1] > ints[i] {
			return false
		}
	}
	return true
}

func isSortedDec(ints []int) bool {
	for i := 1; i < len(ints); i++ {
		if ints[i-1] < ints[i] {
			return false
		}
	}
	return true
}
