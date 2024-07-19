package auth

import (
	"context"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/accesstoken"
)

func GetJwtClaimsFromContext(ctx context.Context) (accesstoken.Claims, error) {
	claims, ok := ctx.Value(accesstoken.ClaimsKey).(accesstoken.Claims)
	if !ok || claims.UserID == 0 {
		return accesstoken.Claims{}, ErrNoClaims
	}

	return claims, nil
}
