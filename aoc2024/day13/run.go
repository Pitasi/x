package main

func Star1(g Game) int {
	return findShortestPath(g.Prize, g.A, g.B, 100)
}

func Star2(g Game) int {
	return findShortestPath(g.Prize, g.A, g.B, 0)
}

func findShortestPath(end, move1, move2 Position, limit int) int {
	a, b := solveLinearEqs(move1.X, move2.X, end.X, move1.Y, move2.Y, end.Y)
	if limit > 0 && (a > limit || b > limit) {
		return 0
	}
	return 3*a + b
}

// aX1 + bX2 = Xn
// aY1 + bY2 = Yn
func solveLinearEqs(x1, x2, xn, y1, y2, yn int) (int, int) {
	det := x1*y2 - x2*y1
	if det == 0 {
		return 0, 0
	}

	aN := (xn*y2 - x2*yn)
	bN := (x1*yn - xn*y1)

	if aN%det != 0 || bN%det != 0 {
		return 0, 0
	}

	a := aN / det
	b := bN / det
	return a, b
}
