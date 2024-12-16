package main

import (
	"fmt"
)

func Run(m Map, robot Position, moves []Move) int {
	for _, move := range moves {
		stepV1(m, &robot, move)
	}

	return scoreV1(m)
}

func RunV2(m Map, robot Position, moves []Move) int {
	for _, move := range moves {
		stepV2(m, &robot, move)
	}

	return scoreV2(m)
}

func stepV1(m Map, robot *Position, move Move) {
	m[robot.Y][robot.X] = '.'
	switch move {
	case Up:
		switch m[robot.Y-1][robot.X] {
		case '.':
			robot.Y--
		case 'O':
			boxX := robot.X
			boxY := robot.Y - 1
			for boxY >= 0 && m[boxY][boxX] == 'O' {
				boxY--
			}
			if boxY >= 0 && m[boxY][boxX] == '.' {
				m[boxY][boxX] = 'O'
				robot.Y--
			}
		case '#':
			// nothing
		default:
			panic(fmt.Errorf("unexpected character '%c'", m[robot.Y-1][robot.X]))
		}
	case Down:
		switch m[robot.Y+1][robot.X] {
		case '.':
			robot.Y++
		case 'O':
			boxX := robot.X
			boxY := robot.Y + 1
			for boxY < len(m) && m[boxY][boxX] == 'O' {
				boxY++
			}
			if boxY < len(m) && m[boxY][boxX] == '.' {
				m[boxY][boxX] = 'O'
				robot.Y++
			}
		case '#':
			// nothing
		default:
			panic(fmt.Errorf("unexpected character '%c'", m[robot.Y+1][robot.X]))
		}
	case Left:
		switch m[robot.Y][robot.X-1] {
		case '.':
			robot.X--
		case 'O':
			boxX := robot.X - 1
			boxY := robot.Y
			for boxX >= 0 && m[boxY][boxX] == 'O' {
				boxX--
			}
			if boxX >= 0 && m[boxY][boxX] == '.' {
				m[boxY][boxX] = 'O'
				robot.X--
			}
		case '#':
			// nothing
		default:
			panic(fmt.Errorf("unexpected character '%c'", m[robot.Y][robot.X-1]))
		}
	case Right:
		switch m[robot.Y][robot.X+1] {
		case '.':
			robot.X++
		case 'O':
			boxX := robot.X + 1
			boxY := robot.Y
			for boxX < len(m[0]) && m[boxY][boxX] == 'O' {
				boxX++
			}
			if boxX < len(m[0]) && m[boxY][boxX] == '.' {
				m[boxY][boxX] = 'O'
				robot.X++
			}
		case '#':
			// nothing
		default:
			panic(fmt.Errorf("unexpected character '%c'", m[robot.Y][robot.X+1]))
		}
	}
	m[robot.Y][robot.X] = '@'
}

func scoreV1(m Map) int {
	var s int
	for y := range m {
		for x := range m[y] {
			if m[y][x] == 'O' {
				s += x + 100*y
			}
		}
	}
	return s
}

func scoreV2(m Map) int {
	var s int
	for y := range m {
		for x := range m[y] {
			if m[y][x] == '[' {
				s += x + 100*y
			}
		}
	}
	return s
}

func stepV2(m Map, robot *Position, move Move) {
	m[robot.Y][robot.X] = '.'
	switch move {
	case Left:
		switch m[robot.Y][robot.X-1] {
		case '.':
			robot.X--
		case ']':
			if canMoveBoxLeft(m, robot.X-2, robot.Y) {
				moveBoxLeft(m, robot.X-2, robot.Y)
				robot.X--
			}
		case '#':
			// nothing
		default:
			panic(fmt.Errorf("unexpected character '%c'", m[robot.Y][robot.X-1]))
		}
	case Right:
		switch m[robot.Y][robot.X+1] {
		case '.':
			robot.X++
		case '[':
			if canMoveBoxRight(m, robot.X+1, robot.Y) {
				moveBoxRight(m, robot.X+1, robot.Y)
				robot.X++
			}
		case '#':
			// nothing
		default:
			panic(fmt.Errorf("unexpected character '%c'", m[robot.Y][robot.X-1]))
		}
	case Up:
		switch m[robot.Y-1][robot.X] {
		case '.':
			robot.Y--
		case '[', ']':
			boxLX := robot.X
			boxLY := robot.Y - 1
			if m[robot.Y-1][robot.X] == ']' {
				boxLX--
			}
			if canMoveBoxVert(m, boxLX, boxLY, -1) {
				moveBoxVert(m, boxLX, boxLY, -1)
				robot.Y--
			}
		}
	case Down:
		switch m[robot.Y+1][robot.X] {
		case '.':
			robot.Y++
		case '[', ']':
			boxLX := robot.X
			boxLY := robot.Y + 1
			if m[robot.Y+1][robot.X] == ']' {
				boxLX--
			}
			if canMoveBoxVert(m, boxLX, boxLY, +1) {
				moveBoxVert(m, boxLX, boxLY, +1)
				robot.Y++
			}
		}
	}
	m[robot.Y][robot.X] = '@'
}

func canMoveBoxLeft(m Map, boxLX, boxY int) bool {
	if m[boxY][boxLX] != '[' {
		panic(fmt.Errorf("expected '[' at %d,%d. Got '%c'", boxLX, boxY, m[boxY][boxLX]))
	}

	// .[]
	if m[boxY][boxLX-1] == '.' {
		return true
	}

	// #[]
	if m[boxY][boxLX-1] == '#' {
		return false
	}

	// [][]
	if m[boxY][boxLX-2] == '[' {
		return canMoveBoxLeft(m, boxLX-2, boxY)
	}

	return true
}

func canMoveBoxRight(m Map, boxLX, boxY int) bool {
	if m[boxY][boxLX] != '[' {
		panic(fmt.Errorf("expected '[' at %d,%d. Got '%c'", boxLX, boxY, m[boxY][boxLX]))
	}

	// [].
	if m[boxY][boxLX+2] == '.' {
		return true
	}

	// []#
	if m[boxY][boxLX+2] == '#' {
		return false
	}

	// [][]
	if m[boxY][boxLX+2] == '[' {
		return canMoveBoxRight(m, boxLX+2, boxY)
	}

	return true
}

func moveBoxLeft(m Map, boxLX, boxY int) {
	if m[boxY][boxLX] != '[' {
		panic(fmt.Errorf("expected '[' at %d,%d. Got '%c'", boxLX, boxY, m[boxY][boxLX]))
	}

	if m[boxY][boxLX-2] == '[' {
		moveBoxLeft(m, boxLX-2, boxY)
	}

	m[boxY][boxLX+1] = '.'
	m[boxY][boxLX] = ']'
	m[boxY][boxLX-1] = '['
}

func moveBoxRight(m Map, boxLX, boxY int) {
	if m[boxY][boxLX] != '[' {
		panic(fmt.Errorf("expected '[' at %d,%d. Got '%c'", boxLX, boxY, m[boxY][boxLX]))
	}

	if m[boxY][boxLX+2] == '[' {
		moveBoxRight(m, boxLX+2, boxY)
	}

	m[boxY][boxLX] = '.'
	m[boxY][boxLX+1] = '['
	m[boxY][boxLX+2] = ']'
}

func canMoveBoxVert(m Map, boxLX, boxY, deltaY int) bool {
	if m[boxY][boxLX] != '[' {
		panic(fmt.Errorf("expected '[' at %d,%d. Got '%c'", boxLX, boxY, m[boxY][boxLX]))
	}
	boxRX := boxLX + 1
	newBoxY := boxY + deltaY

	// #
	if m[newBoxY][boxLX] == '#' || m[newBoxY][boxRX] == '#' {
		return false
	}

	// []
	// []
	if m[newBoxY][boxLX] == '[' {
		return canMoveBoxVert(m, boxLX, newBoxY, deltaY)
	}

	//  []
	//  .[]
	if m[newBoxY][boxRX] == '[' && !canMoveBoxVert(m, boxRX, newBoxY, deltaY) {
		return false
	}

	//  []
	// [].
	if m[newBoxY][boxLX] == ']' && !canMoveBoxVert(m, boxLX-1, newBoxY, deltaY) {
		return false
	}

	return true
}

func moveBoxVert(m Map, boxLX, boxY, deltaY int) {
	if m[boxY][boxLX] != '[' {
		panic(fmt.Errorf("expected '[' at %d,%d, got: '%c'", boxLX, boxY, m[boxY][boxLX]))
	}
	boxRX := boxLX + 1
	newBoxY := boxY + deltaY

	// []
	// ..
	if m[newBoxY][boxLX] == '.' && m[newBoxY][boxRX] == '.' {
		m[newBoxY][boxLX] = '['
		m[newBoxY][boxRX] = ']'
		m[boxY][boxLX] = '.'
		m[boxY][boxRX] = '.'
		return
	}

	// []
	// []
	if m[newBoxY][boxLX] == '[' {
		moveBoxVert(m, boxLX, newBoxY, deltaY)
	}

	//  []
	//  .[]
	if m[newBoxY][boxRX] == '[' {
		moveBoxVert(m, boxRX, newBoxY, deltaY)
	}

	//  []
	// [].
	if m[newBoxY][boxLX] == ']' {
		moveBoxVert(m, boxLX-1, newBoxY, deltaY)
	}

	moveBoxVert(m, boxLX, boxY, deltaY) // try again this box
}
