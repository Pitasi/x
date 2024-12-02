package automap

import (
	"reflect"
)

var mappings map[reflect.Type]map[reflect.Type]any

func Register[A, B any](mapper func(A) B) {
	var (
		a A
		b B
	)
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	if mappings == nil {
		mappings = make(map[reflect.Type]map[reflect.Type]any)
	}
	if mappings[aType] == nil {
		mappings[aType] = make(map[reflect.Type]any)
	}

	mappings[aType][bType] = mapper
}

func Map[A, B any](a A, b *B) {
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(*b)

	f := mappings[aType][bType].(func(A) B)
	*b = f(a)
}
