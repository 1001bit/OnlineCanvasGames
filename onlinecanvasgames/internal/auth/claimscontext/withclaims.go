package claimscontext

import (
	"context"
	"errors"
)

type contextKey string

const (
	usernameKey contextKey = "username"
	userIdKey   contextKey = "userID"
)

func GetContext(ctx context.Context, userID int, username string) context.Context {
	ctx = context.WithValue(ctx, userIdKey, userID)
	ctx = context.WithValue(ctx, usernameKey, username)
	return ctx
}

func GetClaims(ctx context.Context) (int, string, error) {
	userID, ok := ctx.Value(userIdKey).(int)
	if userID == 0 || !ok {
		return 0, "", errors.New("no claims")
	}
	return userID, ctx.Value(usernameKey).(string), nil
}
