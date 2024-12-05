package main

import (
	"reflect"
	"testing"
)

func TestDepsTree(t *testing.T) {
	rules := []PrecedenceRule{
		{47, 53},
		{97, 13},
		{97, 61},
		{97, 47},
		{75, 29},
		{61, 13},
		{75, 53},
		{29, 13},
		{97, 29},
		{53, 29},
		{61, 53},
		{97, 53},
		{61, 29},
		{47, 13},
		{75, 47},
		{97, 75},
		{47, 61},
		{75, 61},
		{47, 29},
		{75, 13},
		{53, 13},
	}

	pageSet := make(map[int]struct{})
	for _, page := range rules {
		pageSet[page.Before] = struct{}{}
		pageSet[page.After] = struct{}{}
	}
	pages := make([]Page, 0, len(pageSet))
	for page := range pageSet {
		pages = append(pages, Page{Number: page})
	}

	tree := BuildDepsTree(rules, pages)

	for _, rule := range rules {
		scoreB := tree.Get(Page{Number: rule.Before})
		scoreA := tree.Get(Page{Number: rule.After})
		if scoreA.IsAncestorOf(scoreB) {
			t.Log(tree)
			t.Fatalf("page %d should come before %d", rule.Before, rule.After)
		}
	}
}

func TestSort(t *testing.T) {
	rules := []PrecedenceRule{
		{47, 53},
		{97, 13},
		{97, 61},
		{97, 47},
		{75, 29},
		{61, 13},
		{75, 53},
		{29, 13},
		{97, 29},
		{53, 29},
		{61, 53},
		{97, 53},
		{61, 29},
		{47, 13},
		{75, 47},
		{97, 75},
		{47, 61},
		{75, 61},
		{47, 29},
		{75, 13},
		{53, 13},
	}
	pj := PrintJob{
		Pages: []Page{
			{Number: 97},
			{Number: 13},
			{Number: 75},
			{Number: 29},
			{Number: 47},
		},
	}

	sortedPj := Sort(rules, pj)
	expected := PrintJob{
		Pages: []Page{
			{Number: 97},
			{Number: 75},
			{Number: 47},
			{Number: 29},
			{Number: 13},
		},
	}

	if !reflect.DeepEqual(sortedPj, expected) {
		t.Fatal("expected", expected, "got", sortedPj)
	}
}
