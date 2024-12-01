package main

import (
	"math/rand"
	"sort"
	"testing"
)

func BenchmarkSort(b *testing.B) {
	rands := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		rands[i] = rand.Intn(10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := make([]int, 0, 10000)
		copy(l, rands)
		sort.Ints(l)
	}
}
