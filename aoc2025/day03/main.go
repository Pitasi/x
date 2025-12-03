package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	input, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	star(input, findMaxJoltage)

	_, _ = input.Seek(0, 0)
	star(input, findMaxJoltageTwelve)
}

func star(input io.Reader, fn func(string) int) {
	scanner := bufio.NewScanner(input)

	var sum int
	for scanner.Scan() {
		line := scanner.Text()
		num := fn(line)
		sum += num
	}

	fmt.Println("star:", sum)
}

func findMaxJoltage(s string) int {
	var (
		c1 rune
		c2 rune

		c1Pos int
	)

	for i, c := range s[:len(s)-1] {
		if c > c1 {
			c1 = c
			c1Pos = i
		}

		if c1 == '9' {
			break
		}
	}

	for _, c := range s[c1Pos+1:] {
		if c > c2 {
			c2 = c
		}
		if c2 == '9' {
			break
		}
	}

	return int(c1-'0')*10 + int(c2-'0')
}

func findMaxJoltageTwelve(s string) int {
	var (
		runes [12]rune
		pos   [12]int
	)

loop:
	for i := range 12 {
		var (
			startScan int
			endScan   = len(s) - 11 + i
		)

		if i != 0 {
			startScan = pos[i-1] + 1
		}

		for cPos, c := range s[startScan:endScan] {
			if c > runes[i] {
				runes[i] = c
				pos[i] = cPos + startScan
			}
			if c == '9' {
				continue loop
			}
		}
	}

	return int(runes[0]-'0')*100000000000 +
		int(runes[1]-'0')*10000000000 +
		int(runes[2]-'0')*1000000000 +
		int(runes[3]-'0')*100000000 +
		int(runes[4]-'0')*10000000 +
		int(runes[5]-'0')*1000000 +
		int(runes[6]-'0')*100000 +
		int(runes[7]-'0')*10000 +
		int(runes[8]-'0')*1000 +
		int(runes[9]-'0')*100 +
		int(runes[10]-'0')*10 +
		int(runes[11]-'0')
}
