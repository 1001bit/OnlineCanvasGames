package auth

import (
	"context"
	"errors"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/accesstoken"
)

var ErrNoClaims = errors.New("no claims found in context")

func GetContextClaims(ctx context.Context) (accesstoken.Claims, error) {
	claims, ok := ctx.Value(accesstoken.ClaimsKey).(accesstoken.Claims)
	if !ok || claims.UserID == 0 {
		return accesstoken.Claims{}, ErrNoClaims
	}

	return claims, nil
}
