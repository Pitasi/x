package main

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	input := []Token{
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

	expected := Program{
		Rules: []PrecedenceRule{
			{47, 53},
			{97, 13},
			{97, 61},
		},
		PrintJobs: []PrintJob{
			{
				Pages: []Page{
					{Number: 75},
					{Number: 47},
					{Number: 61},
					{Number: 53},
					{Number: 29},
				},
			},
			{
				Pages: []Page{
					{Number: 75},
					{Number: 29},
					{Number: 13},
				},
			},
		},
	}

	program, err := Parse(input)
	if err != nil {
		t.Fatal("parse error:", err)
	}

	if !reflect.DeepEqual(program, expected) {
		t.Fatalf("expected %v, got %v", expected, program)
	}
}
