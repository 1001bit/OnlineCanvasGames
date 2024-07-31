package claimscontext

import (
	"context"
)

type contextKey string

const (
	usernameKey contextKey = "username"
)

func GetContext(ctx context.Context, username string) context.Context {
	ctx = context.WithValue(ctx, usernameKey, username)
	return ctx
}

func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}
