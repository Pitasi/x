package main

import "fmt"

func main() {
	t := Tokenizer{input: input}
	var sum int
	for {
		mul, ok := t.Next()
		if !ok {
			break
		}
		sum += execMul(mul)
	}
	fmt.Println(sum)
}

func execMul(m Mul) int {
	return m.A * m.B
}
