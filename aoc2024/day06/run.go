package main

import (
	"errors"
)

func Run(m Map) (int, error) {
	s := simulation{m: m}

	for {
		stop, err := s.step()
		if err != nil {
			return 0, err
		}
		if stop {
			break
		}
	}

	var count int
	for _, row := range m.grid {
		for _, cell := range row {
			if cell, ok := cell.(Marker); ok {
				if cell.IsMarked() {
					count++
				}
			}
		}
	}

	return count, nil
}

func CountLoops(m Map) int {
	// try swapping each cell with a Wall and count loops
	tmpObstacle := &Wall{}
	var loopsCount int
	nRows := len(m.grid)
	nCols := len(m.grid[0])
	for y := range nCols {
		for x := range nRows {
			if x == m.guard.position.x && y == m.guard.position.y {
				continue
			}
			if _, ok := m.grid[y][x].(Obstacle); ok {
				continue
			}
			newMap := m.With(x, y, tmpObstacle)
			_, err := Run(newMap)
			if err == ErrLoopDetected {
				loopsCount++
			}
		}
	}
	return loopsCount
}

type simulation struct {
	m Map
}

var (
	ErrLoopDetected = errors.New("loop detected")
)

func (s *simulation) step() (bool, error) {
	// mark current cell
	currentCell := s.m.GetCell(s.m.guard.position)

	if currentCell == nil {
		return true, nil
	}

	if currentCell, ok := currentCell.(Marker); ok {
		if currentCell.IsMarkedFor(s.m.guard.direction) {
			return false, ErrLoopDetected
		}

		currentCell.Mark(s.m.guard.direction)
	}

	// try step in current direction
	var nextPos Position
	switch s.m.guard.direction {
	case Up:
		nextPos = s.m.guard.position.Up()
	case Down:
		nextPos = s.m.guard.position.Down()
	case Left:
		nextPos = s.m.guard.position.Left()
	case Right:
		nextPos = s.m.guard.position.Right()
	}

	nextCell := s.m.GetCell(nextPos)

	// if next cell is an obstacle, turn 90 degrees
	if IsObstacle(nextCell) {
		switch s.m.guard.direction {
		case Up:
			s.m.guard.SetDirection(Right)
		case Down:
			s.m.guard.SetDirection(Left)
		case Left:
			s.m.guard.SetDirection(Up)
		case Right:
			s.m.guard.SetDirection(Down)
		}
	} else {
		s.m.guard.SetPosition(nextPos)
	}

	return false, nil
}
