package main

import "iter"

func limit2[K, V any](it iter.Seq2[K, V], limit int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var count int
		for k, v := range it {
			if !yield(k, v) {
				return
			}
			count++
			if count == limit {
				return
			}
		}
	}
}
