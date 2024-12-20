package main

import (
	"strings"
)

type trie struct {
	roots map[byte]*node
}

func newTrie() trie {
	return trie{roots: map[byte]*node{}}
}

func (t *trie) add(design Design) {
	root, found := t.roots[design[0]]
	if found {
		root.add(design[1:])
	} else {
		t.roots[design[0]] = &node{
			c:        design[0],
			children: map[byte]*node{},
			final:    len(design) == 1,
		}
		root = t.roots[design[0]]
		root.add(design[1:])
	}
}

func (t trie) String() string {
	var sb strings.Builder
	for _, r := range t.roots {
		sb.WriteString(r.StringIndented(0))
	}
	return sb.String()
}

type node struct {
	c        byte
	children map[byte]*node
	final    bool
}

func (n *node) add(design Design) {
	if len(design) == 0 {
		n.final = true
		return
	}
	child, found := n.children[design[0]]
	if found {
		child.add(design[1:])
	} else {
		n.children[design[0]] = &node{
			c:        design[0],
			children: map[byte]*node{},
		}
		child = n.children[design[0]]
		child.add(design[1:])
	}
}

func (n node) StringIndented(i int) string {
	var sb strings.Builder
	for range i {
		sb.WriteByte(' ')
	}
	sb.WriteByte(n.c)
	if n.final {
		sb.WriteByte('*')
	}
	sb.WriteByte('\n')
	for _, c := range n.children {
		sb.WriteString(c.StringIndented(i + 1))
	}
	return sb.String()
}
