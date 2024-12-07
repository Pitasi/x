package main

import "testing"

func TestCanBeSolved(t *testing.T) {
	res := CanBeSolved(Equation{
		Result: 292,
		Inputs: []int{11, 6, 16, 20},
	})

	if !res {
		t.Errorf("expected true, got false")
	}
}

func TestCanBeSolvedExample(t *testing.T) {
	input := `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`
	equations := parse([]byte(input))

	var sum int
	for equation := range equations {
		if CanBeSolved(equation) {
			sum += equation.Result
		}
	}

	if sum != 3749 {
		t.Errorf("expected 3749, got %d", sum)
	}
}
