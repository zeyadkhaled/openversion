// Package generator provides generators for different purposes.
package generator

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
)

func RandomUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func IsUUID(s string) bool {
	id, err := uuid.Parse(s)
	if err != nil {
		return false
	}

	// Parse allow multiple uuid forms but String function return as
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx which is only allowed form for this
	// project. So compare original with String functions result.
	if s != id.String() {
		return false
	}

	return true
}

func RandomCode() (string, error) {
	var v uint64
	err := binary.Read(rand.Reader, binary.BigEndian, &v)
	if err != nil {
		return "", err
	}

	// since we are using modulus some numbers have more probable.
	return fmt.Sprintf("%05d", int(v%100000)), nil
}

func Fixed(val string) func() (string, error) {
	return func() (string, error) {
		return val, nil
	}
}

func Must(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}
