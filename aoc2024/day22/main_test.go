package main

import "testing"

func TestMix(t *testing.T) {
	seed := 42
	v := 15
	expected := 37
	actual := mix(v, seed)
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func TestRandom(t *testing.T) {
	seed := 123

	expected := []int{
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254,
	}

	for _, expected := range expected {
		actual := random(seed)
		if actual != expected {
			t.Errorf("expected %d, got %d", expected, actual)
		}
		seed = actual
	}
}

func TestNth(t *testing.T) {
	cases := []struct {
		seed     int
		n        int
		expected int
	}{
		{1, 2000, 8685429},
		{10, 2000, 4700978},
		{100, 2000, 15273692},
		{2024, 2000, 8667524},
	}

	for _, c := range cases {
		actual := nth(c.seed, c.n)
		if actual != c.expected {
			t.Errorf("expected %d, got %d", c.expected, actual)
		}
	}
}
