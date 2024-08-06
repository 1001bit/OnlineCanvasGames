package middleware

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/token"
)

func ClaimsContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get access cookie
		cookie, err := r.Cookie("access")

		// if bad cookie
		if err != nil {
			handleTokenRefresh(w, r, next)
			return
		}

		// get claims from cookie
		claims, err := token.ValidateAccessToken(cookie.Value)
		if err != nil {
			handleTokenRefresh(w, r, next)
			return
		}

		ctx := claimscontext.GetContext(r.Context(), claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleTokenRefresh(w http.ResponseWriter, r *http.Request, next http.Handler) {
	username, err := refreshTokens(w, r)
	if err != nil {
		next.ServeHTTP(w, r)
		return
	}

	ctx := claimscontext.GetContext(r.Context(), username)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func refreshTokens(w http.ResponseWriter, r *http.Request) (string, error) {
	// get refresh token
	cookie, err := r.Cookie("refresh")
	if err != nil {
		return "", err
	}

	// refresh tokens
	accessTokenString, refreshTokenString, err := token.RefreshTokens(cookie.Value)
	if err != nil {
		return "", err
	}

	// set new cookies
	accessTokenCookie := token.GenerateAccessTokenCookie(accessTokenString, false)
	refreshTokenCookie := token.GenerateRefreshTokenCookie(refreshTokenString, false)

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	// get claims
	claims, err := token.ValidateAccessToken(accessTokenString)
	if err != nil {
		return "", err
	}

	return claims.Username, nil
}
