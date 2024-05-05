package middleware

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
)

// plain text for unauthorized
func AuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(auth.ClaimsKey).(auth.Claims)
		if !ok {
			api.ServeTextMessage(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// auth page for unauthorized
func AuthHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(auth.ClaimsKey).(auth.Claims)
		if !ok {
			page.HandleAuth(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
