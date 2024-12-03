package main

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		input := "mul(1,2)"
		expected := Mul{1, 2}
		tokenizer := Tokenizer{input: []byte(input)}
		actual, ok := tokenizer.Next()
		if !ok {
			t.Fatal("expected more tokens")
		}
		if actual != expected {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("dont", func(t *testing.T) {
		input := "don't()mul(1,2)"
		tokenizer := Tokenizer{input: []byte(input)}
		_, ok := tokenizer.Next()
		if ok {
			t.Fatal("expected no tokens")
		}
	})

	t.Run("do", func(t *testing.T) {
		input := "do()mul(1,2)"
		expected := Mul{1, 2}
		tokenizer := Tokenizer{input: []byte(input)}
		actual, ok := tokenizer.Next()
		if !ok {
			t.Fatal("expected more tokens")
		}
		if actual != expected {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("do", func(t *testing.T) {
		input := "mul(1,2)don't()mul(3,4)do()mul(5,6)"
		tokenizer := Tokenizer{input: []byte(input)}
		actual, ok := tokenizer.Next()
		if !ok {
			t.Fatal("expected more tokens")
		}
		m := Mul{1, 2}
		if actual != m {
			t.Fatalf("expected %v, got %v", m, actual)
		}

		actual, ok = tokenizer.Next()
		if !ok {
			t.Fatal("expected more tokens")
		}
		m = Mul{5, 6}
		if actual != m {
			t.Fatalf("expected %v, got %v", m, actual)
		}
	})
}
