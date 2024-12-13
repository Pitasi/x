package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	games, err := Parse(input)
	if err != nil {
		panic(err)
	}

	var sum int
	for _, g := range games {
		sum += Star1(g)
	}
	fmt.Println(sum)

	var sum2 int
	games = star2Games(games)
	for _, g := range games {
		sum2 += Star2(g)
	}
	fmt.Println(sum2)
}

func star2Games(in []Game) []Game {
	games := make([]Game, len(in))
	for i, g := range in {
		games[i] = Game{
			A: g.A,
			B: g.B,
			Prize: Position{
				X: g.Prize.X + 10000000000000,
				Y: g.Prize.Y + 10000000000000,
			},
		}
	}
	return games
}
