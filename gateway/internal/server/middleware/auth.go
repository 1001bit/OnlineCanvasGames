package middleware

import (
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/auth/claimscontext"
	"github.com/1001bit/ocg-gateway-service/internal/server/handler/api"
	"github.com/1001bit/ocg-gateway-service/internal/server/handler/page"
)

func AuthHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := claimscontext.GetUsername(r.Context())

		if !ok {
			page.HandleAuth(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := claimscontext.GetUsername(r.Context())

		if !ok {
			api.ServeTextMessage(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
