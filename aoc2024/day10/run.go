package main

type Node struct {
	Value    uint8
	Children []*Node
}

func Roots(grid [][]*Node) []*Node {
	var roots []*Node
	for _, row := range grid {
		for _, node := range row {
			if node.Value == 0 {
				roots = append(roots, node)
			}
		}
	}
	return roots
}

func HikingTrails(head *Node, topHeight uint8) []*Node {
	if head == nil || head.Value != 0 {
		panic("head is nil or head.Value is not 0")
	}

	nextNodes := []*Node{head}
	for i := uint8(1); i < topHeight+1; i++ {
		if len(nextNodes) == 0 {
			break
		}
		var foundNodes []*Node
		for _, n := range nextNodes {
			foundNodes = append(foundNodes, searchChildren(n, i)...)
		}
		nextNodes = foundNodes
	}

	return nextNodes
}

func CountUniques(nodes []*Node) int {
	seen := make(map[*Node]struct{})
	var result int
	for _, n := range nodes {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		result++
	}
	return result
}

func searchChildren(n *Node, searchValue uint8) []*Node {
	if n == nil {
		return nil
	}
	var result []*Node
	for _, child := range n.Children {
		if child.Value == searchValue {
			result = append(result, child)
		}
	}
	return result
}
