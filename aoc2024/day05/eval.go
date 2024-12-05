package main

import (
	"sort"
)

func Sort(rules []PrecedenceRule, pj PrintJob) PrintJob {
	deps := BuildDepsTree(rules, pj.Pages)
	pages := pj.Pages[:]
	sort.Slice(pages, func(i, j int) bool {
		return deps.Get(pages[i]).IsAncestorOf(deps.Get(pages[j]))
	})
	return PrintJob{Pages: pages}
}

func IsSorted(rules []PrecedenceRule, pj PrintJob) bool {
	deps := BuildDepsTree(rules, pj.Pages)
	return sort.SliceIsSorted(pj.Pages, func(i, j int) bool {
		return deps.Get(pj.Pages[i]).IsAncestorOf(deps.Get(pj.Pages[j]))
	})
}

type DepsTree struct {
	nodes map[int]*Node
}

func BuildDepsTree(rules []PrecedenceRule, pages []Page) DepsTree {
	pageSet := make(map[int]struct{})
	for _, page := range pages {
		pageSet[page.Number] = struct{}{}
	}

	nodes := make(map[int]*Node)
	for _, rule := range rules {
		if _, found := pageSet[rule.Before]; !found {
			continue
		}
		if _, found := pageSet[rule.After]; !found {
			continue
		}

		nodeBefore := nodes[rule.Before]
		if nodeBefore == nil {
			// create node for the Before page
			nodeBefore = &Node{Value: rule.Before}
			nodes[rule.Before] = nodeBefore
		}

		nodeAfter := nodes[rule.After]
		if nodeAfter == nil {
			// create node for the After page
			nodeAfter = &Node{Value: rule.After}
			nodes[rule.After] = nodeAfter
		}

		// add the After node as a child of the Before node
		nodeBefore.Add(nodeAfter)
	}

	return DepsTree{
		nodes: nodes,
	}
}

func (p DepsTree) Get(page Page) *Node {
	return p.nodes[page.Number]
}
