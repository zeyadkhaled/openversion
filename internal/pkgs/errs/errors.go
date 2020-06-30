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

type Kind uint8

const (
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
	return k == KindInternal
}
