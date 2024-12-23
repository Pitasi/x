package main

import "testing"

var example = []byte(`kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn`)

func TestExample(t *testing.T) {
	g := Parse(example)
	expected := 7
	actual := star1(g)
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestExample2(t *testing.T) {
	g := Parse(example)
	expected := "co,de,ka,ta"
	actual := star2(g)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
