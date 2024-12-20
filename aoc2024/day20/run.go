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

			if _, ok := candidates[s]; !ok {
				candidates[s] = make(map[Vec]int)
			}
			candidates[s][e] = minCost
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

func CheatV2Saving(g *Grid, c CheatV2) int {
	return abs(g.Cell(c.Start.X, c.Start.Y).remainingTime-g.Cell(c.End.X, c.End.Y).remainingTime) - c.Cost
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
