package main

import (
	"fmt"

	"github.com/casbin/govaluate"
)

var goval = Govaluate{}

type Govaluate struct{}

func (Govaluate) Check(string) bool { return true }

func (Govaluate) Run(input string) error {
	expression, err := govaluate.NewEvaluableExpression(input)
	if err != nil {
		return err
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
