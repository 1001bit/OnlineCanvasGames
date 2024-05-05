package middleware

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
)

func InjectJWTClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.ClaimsFromRequest(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := auth.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
