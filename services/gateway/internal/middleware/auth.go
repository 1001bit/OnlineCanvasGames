package middleware

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/jsonapi"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/xmlapi"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
)

func AuthHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := claimscontext.GetUsername(r.Context())

		if !ok {
			xmlapi.HandleAuthPage(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := claimscontext.GetUsername(r.Context())

		if !ok {
			jsonapi.HandleUnauthorized(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
