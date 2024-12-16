package main

func Step(n int, robot Robot, bounds Vector) Robot {
	return Robot{
		P: robot.P.Add(robot.V.Mul(n)).Mod(bounds),
		V: robot.V,
	}
}

func StepMultiple(steps int, robots []Robot, bounds Vector) []Robot {
	finalRobots := make([]Robot, len(robots))
	for i, robot := range robots {
		finalRobots[i] = Step(steps, robot, bounds)
	}
	return finalRobots
}

func Star1(robots []Robot) int {
	steps := 100
	mapSize := v(101, 103)

	finalRobots := StepMultiple(steps, robots, mapSize)

	var q1, q2, q3, q4 int
	for _, robot := range finalRobots {
		if robot.P.X == mapSize.X/2 || robot.P.Y == mapSize.Y/2 {
			continue
		}

		if robot.P.X < mapSize.X/2 && robot.P.Y < mapSize.Y/2 {
			q1++
		}
		if robot.P.X > mapSize.X/2 && robot.P.Y < mapSize.Y/2 {
			q2++
		}
		if robot.P.X < mapSize.X/2 && robot.P.Y > mapSize.Y/2 {
			q3++
		}
		if robot.P.X > mapSize.X/2 && robot.P.Y > mapSize.Y/2 {
			q4++
		}
	}

	return q1 * q2 * q3 * q4
}
