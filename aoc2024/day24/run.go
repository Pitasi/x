package main

import (
	"fmt"
	"slices"
	"strings"
)

func Star1(g map[string]Gate) uint {
	zGates := FindZGates(g, 'z')
	results := make([]bool, len(zGates))
	for i, name := range zGates {
		results[i] = g[name].Value()
	}
	return Combine(results...)
}

func printsum(g map[string]Gate) {
	xGates := FindZGates(g, 'x')
	yGates := FindZGates(g, 'y')
	zGates := FindZGates(g, 'z')

	xRes := make([]bool, len(xGates))
	for i, name := range xGates {
		xRes[i] = g[name].Value()
	}

	yRes := make([]bool, len(yGates))
	for i, name := range yGates {
		yRes[i] = g[name].Value()
	}

	zRes := make([]bool, len(zGates))
	for i, name := range zGates {
		zRes[i] = g[name].Value()
	}

	x := Combine(xRes...)
	y := Combine(yRes...)
	z := Combine(zRes...)

	fmt.Printf("x: %046b\n", x)
	fmt.Printf("y: %046b\n", y)
	fmt.Printf("c: %046b\n", x+y)
	fmt.Printf("z: %046b\n", z)
}

func sum(g map[string]Gate) (uint, uint, uint) {
	xGates := FindZGates(g, 'x')
	yGates := FindZGates(g, 'y')
	zGates := FindZGates(g, 'z')

	xRes := make([]bool, len(xGates))
	for i, name := range xGates {
		xRes[i] = g[name].Value()
	}

	yRes := make([]bool, len(yGates))
	for i, name := range yGates {
		yRes[i] = g[name].Value()
	}

	zRes := make([]bool, len(zGates))
	for i, name := range zGates {
		zRes[i] = g[name].Value()
	}

	x := Combine(xRes...)
	y := Combine(yRes...)
	z := Combine(zRes...)

	return x, y, z
}

func Star2(g map[string]Gate) string {
	// most of the work done outside of the code:
	// I took the input, draw it using graphviz, and then
	// visually checked the wrong gates:
	//  - z gates should be XOR

	// Also confirmed by checking the expected output:
	// printsum(g)
	//                34                        12     8 7
	// x: 01011110100  1  101001111010101111110  1 111 0 1 1001101
	// y: 01101111011  1  000111100111100111101  0 100 0 0 0011001
	// c: 11001110000  0  110001100010010111100  0 011 0 1 1100110
	// z: 11001101111  1  110001100010010111011  1 011 1 0 1100110
	//
	// From the above we can pinpoint the wrong bits and go up from them.

	swaps := [][]string{
		{"z07", "rts"},
		{"z12", "jpj"},
		{"kgj", "z26"},
		{"vvw", "chv"},
	}
	for _, s := range swaps {
		swap(g, s[0], s[1])
	}

	mask := uint(0b11111111)
	x, y, z := sum(g)
	if (x+y)&mask != z&mask {
		printsum(g)
		return "error"
	}

	var swapped []string
	for _, s := range swaps {
		swapped = append(swapped, s[0], s[1])
	}
	slices.Sort(swapped)
	return strings.Join(swapped, ",")
}

func swap(g map[string]Gate, n1, n2 string) {
	g[n1], g[n2] = g[n2], g[n1]
}

func deps(g map[string]Gate, gate Gate) []string {
	var d []string
	switch gate := gate.(type) {
	case *GateAnd:
		d = append(d, name(g, gate))
		d = append(d, deps(g, gate.left)...)
		d = append(d, deps(g, gate.right)...)
	case *GateOr:
		d = append(d, name(g, gate))
		d = append(d, deps(g, gate.left)...)
		d = append(d, deps(g, gate.right)...)
	case *GateXor:
		d = append(d, name(g, gate))
		d = append(d, deps(g, gate.left)...)
		d = append(d, deps(g, gate.right)...)
		return d
	case *LazyGate:
		d = append(d, gate.name)
		d = append(d, deps(g, g[gate.name])...)
	}
	dWithoutZ := make([]string, 0, len(d))
	for _, name := range d {
		if name[0] == 'z' || name[0] == 'x' || name[0] == 'y' {
			continue
		}
		dWithoutZ = append(dWithoutZ, name)
	}
	return dWithoutZ
}

func name(g map[string]Gate, gate Gate) string {
	if gate == nil {
		panic("can't be nil")
	}

	for name := range g {
		if g[name] == gate {
			return name
		}
	}

	panic(fmt.Sprintf("gate not found: %T", gate))
}

func FindZGates(g map[string]Gate, prefix byte) []string {
	var gates []string
	for name := range g {
		if name[0] == prefix {
			gates = append(gates, name)
		}
	}
	slices.Sort(gates)
	slices.Reverse(gates)
	return gates
}

func Combine(b ...bool) uint {
	var n uint
	for _, b := range b {
		n <<= 1
		if b {
			n |= 1
		}
	}
	return n
}
