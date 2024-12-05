package main

import "testing"

func TestTokenizer(t *testing.T) {
	input := `47|53
97|13
97|61

75,47,61,53,29
75,29,13`

	expected := []Token{
		{Type: "int", Value: []byte("47")},
		{Type: "pipe", Value: []byte("|")},
		{Type: "int", Value: []byte("53")},
		{Type: "newline", Value: []byte("\n")},
		{Type: "int", Value: []byte("97")},
		{Type: "pipe", Value: []byte("|")},
		{Type: "int", Value: []byte("13")},
		{Type: "newline", Value: []byte("\n")},
		{Type: "int", Value: []byte("97")},
		{Type: "pipe", Value: []byte("|")},
		{Type: "int", Value: []byte("61")},
		{Type: "newline", Value: []byte("\n")},
		{Type: "newline", Value: []byte("\n")},
		{Type: "int", Value: []byte("75")},
		{Type: "comma", Value: []byte(",")},
		{Type: "int", Value: []byte("47")},
		{Type: "comma", Value: []byte(",")},
		{Type: "int", Value: []byte("61")},
		{Type: "comma", Value: []byte(",")},
		{Type: "int", Value: []byte("53")},
		{Type: "comma", Value: []byte(",")},
		{Type: "int", Value: []byte("29")},
		{Type: "newline", Value: []byte("\n")},
		{Type: "int", Value: []byte("75")},
		{Type: "comma", Value: []byte(",")},
		{Type: "int", Value: []byte("29")},
		{Type: "comma", Value: []byte(",")},
		{Type: "int", Value: []byte("13")},
	}

	tokens, err := Tokenize([]byte(input))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(tokens) != len(expected) {
		t.Errorf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type != expected[i].Type {
			t.Errorf("expected token %d to be of type %s, got %s", i, expected[i].Type, tokens[i].Type)
		}

		if string(tokens[i].Value) != string(expected[i].Value) {
			t.Errorf("expected token %d to be %s, got %s", i, string(expected[i].Value), string(tokens[i].Value))
		}
	}
}
