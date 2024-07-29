package middleware

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/claimscontext"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
)

func AuthHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := claimscontext.GetClaims(r.Context())

		if err != nil {
			page.HandleAuth(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := claimscontext.GetClaims(r.Context())

		if err != nil {
			api.ServeTextMessage(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
