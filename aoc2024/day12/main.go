package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	g := Parse(input)
	price := RunStar1(g)
	fmt.Println(price)

	price = RunStar2(g)
	fmt.Println(price)
}
