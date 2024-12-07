package main

import (
	_ "embed"
	"fmt"
	"iter"
	"strconv"
)

//go:embed input.txt
var input []byte

func main() {
	equations := parse(input)

	var sum int
	for equation := range equations {
		if CanBeSolved(equation) {
			sum += equation.Result
		}
	}
	fmt.Println(sum)
}

func CanBeSolved(eq Equation) bool {
	expected := eq.Result

	for op := range AllOperators(len(eq.Inputs) - 1) {
		res := solve(eq.Inputs, op)
		if expected == res {
			return true
		}
	}

	return false
}

func solve(inputs []int, ops Operators) int {
	result := inputs[0]
	for i := 1; i < len(inputs); i++ {
		op := ops.Get(i - 1)
		switch op {
		case 0b00:
			result += inputs[i]
		case 0b01:
			result *= inputs[i]
		case 0b10:
			r, err := strconv.Atoi(strconv.Itoa(result) + strconv.Itoa(inputs[i]))
			if err != nil {
				panic(err)
			}
			result = r
		case 0b11:
			// unused
			continue
		default:
			panic(fmt.Sprintf("invalid op: %b", op))
		}
	}
	return result
}

type Operators uint64

// AllOperators returns a sequence of all possible "count"
// operators.
func AllOperators(count int) iter.Seq[Operators] {
	return func(yield func(Operators) bool) {
		for i := range maxOperatorsValue(count) {
			if !i.IsValid() {
				continue
			}
			if !yield(i) {
				return
			}
		}
	}
}

func maxOperatorsValue(count int) Operators {
	if count > 32 {
		panic("too many inputs for a little uint :(")
	}

	var max Operators
	for i := range count {
		max |= 0b11 << (2 * i)
	}
	return max + 1
}

func (o Operators) IsValid() bool {
	for i := range 32 {
		op := (o & (0b11 << (2 * (i)))) >> (2 * (i))
		if op == 0b11 {
			return false
		}
	}
	return true
}

func (o Operators) Get(idx int) uint8 {
	op := (o & (0b11 << (2 * idx))) >> (2 * idx)
	return uint8(op)
}
