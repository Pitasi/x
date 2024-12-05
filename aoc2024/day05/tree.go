package main

type Node struct {
	Value    int
	Children []*Node
}

func (n *Node) Add(child *Node) {
	n.Children = append(n.Children, child)
}

func (n *Node) IsAncestorOf(maybeDescendant *Node) bool {
	if n == maybeDescendant {
		return false
	}
	for _, child := range n.Children {
		if child == maybeDescendant {
			return true
		}
		if child.IsAncestorOf(maybeDescendant) {
			return true
		}
	}
	return false
}
