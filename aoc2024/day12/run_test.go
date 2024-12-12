package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	input := `AAAA
BBCD
BBCC
EEEC`

	g := Parse([]byte(input))
	price := RunStar1(g)

	expected := 140
	if price != expected {
		t.Errorf("expected %d, got %d", expected, price)
	}
}

func TestRun2(t *testing.T) {
	input := `OOOOO
OXOXO
OOOOO
OXOXO
OOOOO`

	g := Parse([]byte(input))
	price := RunStar1(g)

	expected := 772
	if price != expected {
		t.Errorf("expected %d, got %d", expected, price)
	}
}

func TestComputeAreaSize(t *testing.T) {
	cases := []struct {
		input string
		value byte
		want  int
	}{
		{
			input: `AAAA
BBCD
BBCC
EEEC`,
			value: 'A',
			want:  4,
		},
		{
			input: `AAAA
BBCD
BBCC
EEEC`,
			value: 'B',
			want:  4,
		},
		{
			input: `AAAA
BBCD
BBCC
EEEC`,
			value: 'C',
			want:  8,
		},
		{
			input: `EEEEE
EXXXX
EEEEE
EXXXX
EEEEE`,
			value: 'E',
			want:  12,
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			g := Parse([]byte(tc.input))
			areas := findAreasForValue(g, tc.value)
			var got int
			for _, area := range areas {
				got += computeAreaSides(area)
			}
			if got != tc.want {
				t.Errorf("expected %d, got %d", tc.want, got)
			}
		})
	}
}
