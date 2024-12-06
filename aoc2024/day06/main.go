package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	m := Parse(input)

	c, _ := Run(m)
	fmt.Println(c)

	loops := CountLoops(m)
	fmt.Println(loops)
}
