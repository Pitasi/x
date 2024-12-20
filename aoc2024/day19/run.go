package main

var pathExistsMemo = make(map[string]int)

func PathExists(t trie, design Design) int {
	cacheK := string(design)
	if res, found := pathExistsMemo[cacheK]; found {
		return res
	}

	if len(design) == 0 {
		pathExistsMemo[cacheK] = 1
		return 1
	}

	r, found := t.roots[design[0]]
	if !found {
		pathExistsMemo[cacheK] = 0
		return 0
	}

	submatches := matches(r, design, 1)

	var totalCount int
	for _, m := range submatches {
		unmatched := Design(design[m.count:])
		totalCount += PathExists(t, unmatched)
	}

	pathExistsMemo[cacheK] = totalCount
	return totalCount
}

type match struct {
	count    int
	complete bool
}

func matches(n *node, design Design, count int) []match {
	if len(design) == 0 || len(design) == 1 && n.final {
		return []match{{count, true}}
	}

	if n.c != design[0] {
		panic("expected " + string(n.c) + " got " + string(design[0]))
	}

	if len(design) == 1 {
		// since n is not final can't continue, must stop here
		return nil
	}

	child, found := n.children[design[1]]
	if !found && !n.final {
		// can't stop at node n because it's not final
		// can't continue because can't match the rest
		return nil
	}

	if found && !n.final {
		// can't stop at node 'n' because it's not final
		// must continue to try and match the rest
		return matches(child, design[1:], count+1)
	}

	if !found && n.final {
		// partial match and no other possibilities, stop here
		return []match{{count, false}}
	}

	// n is a final node, submatch exists
	// we could stop here or continue, must try both
	return append(matches(child, design[1:], count+1), match{count, false})
}
