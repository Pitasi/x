package main

import "testing"

func TestCombine(t *testing.T) {
	cases := []struct {
		input []bool
		want  uint
	}{
		{[]bool{false, false, false, false}, 0},
		{[]bool{false, false, false, true}, 1},
		{[]bool{false, false, true, false}, 2},
		{[]bool{false, false, true, true}, 3},
		{[]bool{false, true, false, false}, 4},
		{[]bool{false, true, false, true}, 5},
		{[]bool{false, true, true, false}, 6},
		{[]bool{false, true, true, true}, 7},
		{[]bool{true, false, false, false}, 8},
		{[]bool{true, false, false, true}, 9},
		{[]bool{true, false, true, false}, 10},
		{[]bool{true, false, true, true}, 11},
		{[]bool{true, true, false, false}, 12},
		{[]bool{true, true, false, true}, 13},
		{[]bool{true, true, true, false}, 14},
		{[]bool{true, true, true, true}, 15},
	}

	for _, c := range cases {
		got := Combine(c.input...)
		if got != c.want {
			t.Errorf("expected %d, got %d", c.want, got)
		}
	}
}

func TestExampleSmall(t *testing.T) {
	input := []byte(`x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02`)

	g := Parse(input)
	expected := uint(4)
	actual := Star1(g)
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestExampleLarge(t *testing.T) {
	input := []byte(`x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1

ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk OR frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj`)

	g := Parse(input)
	expected := uint(2024)
	actual := Star1(g)
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}
