package middleware

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	"github.com/1001bit/OnlineCanvasGames/internal/auth/accesstoken"
	"github.com/golang-jwt/jwt/v5"
)

func InjectJWTClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := accesstoken.ClaimsFromRequest(r)

		switch err {
		case nil:
			// no error
		case jwt.ErrTokenExpired, http.ErrNoCookie:
			// Token expired
			claims, err = auth.RefreshTokens(w, r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

		default:
			// Other
			next.ServeHTTP(w, r)
			return
		}

		ctx := accesstoken.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
