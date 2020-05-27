// Package tristate provides a bool with 3 states.
package tristate

import "strconv"

type Bool uint8

const (
	NotSet Bool = iota
	True
	False
)

func NewFromStr(s string) Bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return NotSet
	}
	if b {
		return True
	}
	return False
}

func New(b *bool) Bool {
	if b == nil {
		return NotSet
	}

	if *b {
		return True
	}

	return False
}
