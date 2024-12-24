package main

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02`)

	expected := map[string]Gate{
		"x00": &GateTrue{},
		"x01": &GateTrue{},
		"x02": &GateTrue{},
		"y00": &GateFalse{},
		"y01": &GateTrue{},
		"y02": &GateFalse{},
	}
	expected["z00"] = &GateAnd{left: expected["x00"], right: expected["y00"]}
	expected["z01"] = &GateXor{left: expected["x01"], right: expected["y01"]}
	expected["z02"] = &GateOr{left: expected["x02"], right: expected["y02"]}

	actual := Parse(input)
	if len(actual) != len(expected) {
		t.Errorf("expected %d gates, got %d", len(expected), len(actual))
	}
	for name, expectedGate := range expected {
		actualGate := actual[name]
		if actualGate.Value() != expectedGate.Value() {
			t.Errorf("expected %v, got %v", expectedGate, actualGate)
		}
	}

	z00 := actual["z00"].Value()
	z01 := actual["z01"].Value()
	z02 := actual["z02"].Value()

	expectedZ00 := false
	expectedZ01 := false
	expectedZ02 := true

	if z00 != expectedZ00 {
		t.Errorf("expected z00=%v, got %v", expectedZ00, z00)
	}
	if z01 != expectedZ01 {
		t.Errorf("expected z01=%v, got %v", expectedZ01, z01)
	}
	if z02 != expectedZ02 {
		t.Errorf("expected z02=%v, got %v", expectedZ02, z02)
	}
}
