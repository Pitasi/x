package main

import "fmt"

func MapV2(m Map) (Map, Position) {
	var robot Position
	newMap := make(Map, len(m))
	for y := range m {
		newMap[y] = make([]byte, 2*len(m[y]))
		for x := range m[y] {
			newX := x * 2
			switch m[y][x] {
			case '.':
				newMap[y][newX] = '.'
				newMap[y][newX+1] = '.'
			case 'O':
				newMap[y][newX] = '['
				newMap[y][newX+1] = ']'
			case '#':
				newMap[y][newX] = '#'
				newMap[y][newX+1] = '#'
			case '@':
				robot.X = newX
				robot.Y = y
				newMap[y][newX] = '@'
				newMap[y][newX+1] = '.'
			default:
				panic(fmt.Errorf("unexpected character '%c'", m[y][x]))
			}
		}
	}

	return newMap, robot
}
