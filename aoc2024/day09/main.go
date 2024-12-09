package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	disk := Parse(input)
	CompactByFile(disk)
	fmt.Println(disk.checksum())
}
