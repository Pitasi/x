package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	m := Parse(input)
	count := Run(m)
	fmt.Println(count)
}

func Run(m Map) int {
	// group by antenna frequency
	antennas := make(map[byte][]Position)
	for y, row := range m.grid {
		for x, cell := range row {
			if cell.antenna == 0 {
				continue
			}
			antennas[cell.antenna] = append(antennas[cell.antenna], Position{x: x, y: y})
		}
	}

	// mark antinodes
	nRows := len(m.grid)
	nCols := len(m.grid[0])
	for _, cells := range antennas {
		for _, a := range cells {
			for _, b := range cells {
				if a == b {
					continue
				}

				antinodes1 := a.Antinode(b, nRows, nCols)
				antinodes2 := b.Antinode(a, nRows, nCols)

				for _, node := range antinodes1 {
					cell, ok := m.GetCell(node)
					if ok {
						cell.MarkAntinode()
					}
				}

				for _, node := range antinodes2 {
					cell, ok := m.GetCell(node)
					if ok {
						cell.MarkAntinode()
					}
				}
			}
		}
	}

	// count antinodes
	var count int
	for y := range nCols {
		for x := range nRows {
			cell, ok := m.GetCell(Position{x: x, y: y})
			if ok && cell.antinode {
				count++
			}
		}
	}

	return count
}

func printMap(m Map) {
	for _, row := range m.grid {
		for _, cell := range row {
			if cell.antenna == 0 {
				if cell.antinode {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			} else {
				if cell.antinode {
					fmt.Printf("*")
				} else {
					fmt.Printf("%c", cell.antenna)
				}
			}
		}
		fmt.Println()
	}
}
