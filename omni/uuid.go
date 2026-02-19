package main

import (
	"fmt"

	googleuuid "github.com/google/uuid"
)

var uuid = UUID{}

type UUID struct{}

func (UUID) Check(string) bool { return false }

func (UUID) Run(input string) error {
	fmt.Println(googleuuid.NewString())
	return nil
}
