package main

import (
	"bytes"
	"testing"
)

func TestExample1(t *testing.T) {
	m := M{}
	m.SetRegs(0, 0, 9)
	m.Load([]byte{2, 6})

	m.Run()

	a, b, c := 0, 1, 9
	if m.A != a {
		t.Errorf("expected A=%d, got %d", a, m.A)
	}
	if m.B != b {
		t.Errorf("expected B=%d, got %d", b, m.B)
	}
	if m.C != c {
		t.Errorf("expected C=%d, got %d", c, m.C)
	}
}

func TestExample2(t *testing.T) {
	var out bytes.Buffer
	m := M{Out: &out}
	m.SetRegs(10, 0, 0)
	m.Load([]byte{5, 0, 5, 1, 5, 4})
	m.Run()

	outBytes := out.Bytes()
	expected := []byte{0, 1, 2}
	if !bytes.Equal(outBytes, expected) {
		t.Errorf("expected %v, got %v", expected, outBytes)
	}
}

func TestExample3(t *testing.T) {
	var out bytes.Buffer
	m := M{Out: &out}
	m.SetRegs(2024, 0, 0)
	m.Load([]byte{0, 1, 5, 4, 3, 0})
	m.Run()

	outBytes := out.Bytes()
	expected := []byte{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}
	if !bytes.Equal(outBytes, expected) {
		t.Errorf("expected %v, got %v", expected, outBytes)
	}

	a := 0
	if m.A != a {
		t.Errorf("expected A=%d, got %d", a, m.A)
	}
}

func TestExample4(t *testing.T) {
	var out bytes.Buffer
	m := M{Out: &out}
	m.SetRegs(0, 29, 0)
	m.Load([]byte{1, 7})
	m.Run()

	b := 26
	if m.B != b {
		t.Errorf("expected B=%d, got %d", b, m.B)
	}
}

func TestExample5(t *testing.T) {
	var out bytes.Buffer
	m := M{Out: &out}
	m.SetRegs(0, 2024, 43690)
	m.Load([]byte{4, 0})
	m.Run()

	b := 44354
	if m.B != b {
		t.Errorf("expected B=%d, got %d", b, m.B)
	}
}

func TestExample6(t *testing.T) {
	input := []byte(`Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`)

	a, b, c, program := Parse(input)

	var out bytes.Buffer
	m := M{Out: &out}
	m.SetRegs(a, b, c)
	m.Load(program)
	m.Run()

	outBytes := out.Bytes()
	expected := []byte{4, 6, 3, 5, 6, 3, 5, 2, 1, 0}
	if !bytes.Equal(outBytes, expected) {
		t.Errorf("expected %v, got %v", expected, outBytes)
	}
}

func TestWriter(t *testing.T) {
	writer := &w{expected: []byte{4, 6, 3, 5, 6, 3, 5, 2, 1, 0}}

	input := []byte(`Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`)

	a, b, c, program := Parse(input)

	m := M{Out: writer}
	m.SetRegs(a, b, c)
	m.Load(program)
	m.Run()

	if !writer.Done() {
		t.Errorf("expected writer to be done")
	}
}

func TestRun(t *testing.T) {
	_, b, c, program := Parse(input)

	var out bytes.Buffer
	m := M{Out: &out}
	m.SetRegs(109019476330651, b, c)
	m.Load(program)
	m.Run()

	outBytes := out.Bytes()
	printOut(outBytes)
}
