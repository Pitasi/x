package main

import (
	"bytes"
	"fmt"

	"github.com/cometbft/cometbft/libs/strings"
)

// Tx is what the users will submit. This of this as as the payload of an API
// call.
// In this example, we only have one request type (set a key/value pair).
type Tx struct {
	Key []byte
	Val []byte
}

var ErrInvalidTx = fmt.Errorf("invalid tx format")

// UnmarshalTx unmarshals a Tx from a byte slice.
// Format: key=value
func UnmarshalTx(data []byte) (*Tx, error) {
	parts := bytes.SplitN(data, []byte("="), 2)
	if len(parts) != 2 {
		return nil, ErrInvalidTx
	}
	return &Tx{
		Key: parts[0],
		Val: parts[1],
	}, nil
}

// Valid returns 0 if the Tx is valid.
// A valid tx must have a key of at least 3 characters.
func (tx *Tx) Valid() error {
	if len(tx.Key) < 3 {
		return fmt.Errorf("key too short, must be at least 3 characters")
	}
	if len(tx.Val) == 0 {
		return fmt.Errorf("value too short, must be at least 1 character")
	}
	if !strings.IsASCIIText(string(tx.Key)) {
		return fmt.Errorf("key is not ASCII")
	}
	return nil
}
