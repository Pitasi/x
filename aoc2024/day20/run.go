package main

func CheatV1Candidates(g *Grid) ([]Vec, []Vec) {
	var (
		hCandidates []Vec
		vCandidates []Vec
	)
	for x := 0; x < g.Width(); x++ {
		for y := 0; y < g.Height(); y++ {
			if isWallSeparatingPathH(g, Vec{x, y}) {
				hCandidates = append(hCandidates, Vec{x, y})
			} else if isWallSeparatingPathV(g, Vec{x, y}) {
				vCandidates = append(vCandidates, Vec{x, y})
			}
		}
	}

	return hCandidates, vCandidates
}

func isWallSeparatingPathH(g *Grid, v Vec) bool {
	return g.Cell(v.X, v.Y).IsWall() && !g.Cell(v.X+1, v.Y).IsWall() && !g.Cell(v.X-1, v.Y).IsWall()
}

func isWallSeparatingPathV(g *Grid, v Vec) bool {
	return g.Cell(v.X, v.Y).IsWall() && !g.Cell(v.X, v.Y-1).IsWall() && !g.Cell(v.X, v.Y+1).IsWall()
}

func CheatV1SavingH(g *Grid, wall Vec) int {
	if !g.Cell(wall.X, wall.Y).IsWall() {
		panic("wall is not a wall")
	}
	return saving(g, Vec{wall.X - 1, wall.Y}, Vec{wall.X + 1, wall.Y})
}

func CheatV1SavingV(g *Grid, wall Vec) int {
	if !g.Cell(wall.X, wall.Y).IsWall() {
		panic("wall is not a wall")
	}
	return saving(g, Vec{wall.X, wall.Y - 1}, Vec{wall.X, wall.Y + 1})
}

type CheatV2 struct {
	Start, End Vec
	Cost       int
}

func CheatV2Candidates(g *Grid, minSaving, maxCheatTime int) []CheatV2 {
	candidates := make(map[Vec]map[Vec]int)

	for startIdx := 0; startIdx < len(g.track)-2; startIdx++ {
		for endIdx := startIdx + 1; endIdx < len(g.track); endIdx++ {
			s := g.track[startIdx]
			e := g.track[endIdx]

			minCost := distance(s, e)
			if minCost > maxCheatTime {
				continue
			}

			maxSaving := CheatV2Saving(g, CheatV2{s, e, minCost})
			if maxSaving < minSaving {
				continue
			}

			cost, found := getCheatCost(g, s, e, maxCheatTime)
			if !found {
				continue
			}

			if _, ok := candidates[s]; !ok {
				candidates[s] = make(map[Vec]int)
			}
			candidates[s][e] = cost
		}
	}

	c := make([]CheatV2, 0, len(candidates))
	for s, es := range candidates {
		for e, cost := range es {
			saving := CheatV2Saving(g, CheatV2{s, e, cost})
			if saving < minSaving {
				continue
			}
			c = append(c, CheatV2{s, e, cost})
		}
	}
	return c
}

func getCheatCost(g *Grid, s, e Vec, maxCheatTime int) (int, bool) {
	if g.Cell(s.X, s.Y).remainingTime >= g.Cell(e.X, e.Y).remainingTime {
		// going backwards is not allowed
		return 0, false
	}
	d := distance(s, e)
	if d <= 1 || d > maxCheatTime {
		return 0, false
	}

	return findWallPath(g, s, e, make(map[Vec]int), maxCheatTime)
}

func CheatV2Saving(g *Grid, c CheatV2) int {
	return abs(g.Cell(c.Start.X, c.Start.Y).remainingTime-g.Cell(c.End.X, c.End.Y).remainingTime) - c.Cost
}

func findWallPath(g *Grid, s, e Vec, visited map[Vec]int, availSteps int) (int, bool) {
	if s == e {
		return 0, true
	}

	if availSteps < 0 {
		return 0, false
	}

	visited[s] = availSteps
	moves := []Vec{
		{s.X, s.Y - 1},
		{s.X, s.Y + 1},
		{s.X - 1, s.Y},
		{s.X + 1, s.Y},
	}

	var (
		minCost = 9999999999999
		found   bool
	)
	for _, m := range moves {
		if m == e {
			return 1, true
		}
		if m.X < 0 || m.X >= g.Width() || m.Y < 0 || m.Y >= g.Height() {
			continue
		}
		if !g.Cell(m.X, m.Y).IsWall() {
			continue
		}
		if s, ok := visited[m]; ok && availSteps < s {
			// we already have a better path to m
			continue
		}
		cost, f := findWallPath(g, m, e, visited, availSteps-1)
		cost = cost + 1
		if f && cost < minCost {
			minCost = cost
			found = true
		}
	}

	if !found {
		return 0, false
	}
	return minCost, true
}

func saving(g *Grid, v1, v2 Vec) int {
	c1 := g.Cell(v1.X, v1.Y)
	c2 := g.Cell(v2.X, v2.Y)
	return abs(c1.remainingTime-c2.remainingTime) - distance(v1, v2)
}

func distance(v1, v2 Vec) int {
	return abs(v1.X-v2.X) + abs(v1.Y-v2.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
