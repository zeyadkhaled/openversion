package auth

import (
	"context"
)

type contextKey string

const authContextKey contextKey = "user"

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, authContextKey, user)
}

func UserFromContext(ctx context.Context) User {
	uRaw := ctx.Value(authContextKey)
	if uRaw == nil {
		return User{
			Source: SourceAnonymous,
			ID:     "anonymous-id",
			Phone:  "+908508855647",
			Name:   "Anonymous User",
		}
	}
	return uRaw.(User)
}
