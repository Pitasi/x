package main

import (
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
)

var hex2Dec = HexToDec{}

type HexToDec struct{}

func (HexToDec) Check(input string) bool {
	return strings.HasPrefix(input, "0x")
}

func (HexToDec) Run(input string) error {
	input = strings.TrimPrefix(input, "0x")

	if len(input)%2 != 0 {
		input = "0" + input
	}

	res, err := hex.DecodeString(input)
	if err != nil {
		return err
	}

	slices.Reverse(res)
	var n uint64
	for i, x := range res {
		n += uint64(x) * uint64(pow16(i))
	}

	fmt.Println(n)

	return nil
}

func pow16(n int) int {
	if n == 0 {
		return 1
	}
	return 2 << (4 * n)
}
