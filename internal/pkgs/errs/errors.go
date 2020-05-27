// Package errs error type used in all packages
package errs

import (
	"errors"
	"fmt"
)

type E struct {
	Kind       Kind
	PublicMsg  string
	Parameters []string
	Wrapped    error
}

func (e E) Error() string {
	return fmt.Sprintf("[%s] %v", e.Kind.String(), e.Wrapped)
}

func (e E) Unwrap() error {
	return e.Wrapped
}

// Kind shows kind of error
//
//go:generate stringer -type=Kind
type Kind uint8

const (
	// KindInternal must be 0 for default E to be Internal Error
	// Note: Do not forget to update errshttp with new values
	KindInternal Kind = iota
	KindUnauthorized
	KindLoginPassword
	KindForbidden
	KindNotFound
	KindConflict
	KindDuplicate
	KindParameterErr
	KindTooManyRequests
	KindDependentService
	KindBlocked
)

func Is(err error, k Kind) bool {
	if err == nil {
		return false
	}

	var e E
	if errors.As(err, &e) {
		return e.Kind == k
	}
	// every unknown err is assumed to be internal err
	return k == KindInternal
}
