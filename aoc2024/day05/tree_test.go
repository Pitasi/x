package main

import "testing"

func TestIsAncestorOf(t *testing.T) {
	root := &Node{Value: 1}
	child1 := &Node{Value: 2}
	child2 := &Node{Value: 3}
	child3 := &Node{Value: 4}

	root.Add(child1)
	root.Add(child2)
	child1.Add(child2)
	child1.Add(child3)

	if !root.IsAncestorOf(child1) {
		t.Errorf("expected root to be ancestor of child1")
	}

	if !root.IsAncestorOf(child2) {
		t.Errorf("expected root to be ancestor of child2")
	}

	if !root.IsAncestorOf(child3) {
		t.Errorf("expected root to be ancestor of child3")
	}
}
