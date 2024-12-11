package main

import (
	"fmt"
	"strconv"
)

func main() {
	input := []int{70949, 6183, 4, 3825336, 613971, 0, 15, 182}

	var count1 int
	for _, n := range input {
		count1 += compute(n, 25)
	}
	fmt.Println(count1)

	var count int
	for _, n := range input {
		count += compute(n, 75)
	}
	fmt.Println(count)
}

type pair struct {
	a int
	b int
}

var memo = make(map[pair]int)

func compute(n, count int) int {
	if res, ok := memo[pair{n, count}]; ok {
		return res
	}

	if count == 0 {
		return 1
	}

	if n == 0 {
		res := compute(1, count-1)
		memo[pair{n, count}] = res
		return res
	}

	valueStr := strconv.Itoa(n)
	if len(valueStr)%2 == 0 {
		lHalf := valueStr[:len(valueStr)/2]
		rHalf := valueStr[len(valueStr)/2:]

		lHalfInt, err := strconv.Atoi(lHalf)
		if err != nil {
			panic(err)
		}
		rHalfInt, err := strconv.Atoi(rHalf)
		if err != nil {
			panic(err)
		}

		res := compute(lHalfInt, count-1) + compute(rHalfInt, count-1)
		memo[pair{n, count}] = res
		return res
	}

	res := compute(n*2024, count-1)
	memo[pair{n, count}] = res
	return res
}
