package middleware

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/api"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/components"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
)

func AuthHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := claimscontext.GetUsername(r.Context())

		if !ok {
			components.HandleAuth(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := claimscontext.GetUsername(r.Context())

		if !ok {
			api.HandleUnauthorized(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
