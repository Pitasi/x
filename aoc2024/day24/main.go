package main

import (
	_ "embed"
	"fmt"
)

type Gate interface {
	Value() bool
}

type GateTrue struct{}

func (g *GateTrue) Value() bool { return true }

type GateFalse struct{}

func (g *GateFalse) Value() bool { return false }

type GateAnd struct {
	left, right Gate
}

func (g *GateAnd) Value() bool {
	return g.left.Value() && g.right.Value()
}

type GateOr struct {
	left, right Gate
}

func (g *GateOr) Value() bool {
	return g.left.Value() || g.right.Value()
}

type GateXor struct {
	left, right Gate
}

func (g *GateXor) Value() bool {
	return g.left.Value() != g.right.Value()
}

type LazyGate struct {
	m    map[string]Gate
	name string
}

func (g *LazyGate) Value() bool {
	return g.m[g.name].Value()
}

//go:embed input.txt
var input []byte

func main() {
	g := Parse(input)
	fmt.Println(Star2(g))
}
