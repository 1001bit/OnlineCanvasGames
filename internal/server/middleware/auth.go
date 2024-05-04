package middleware

import (
	"context"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	usermodel "github.com/1001bit/OnlineCanvasGames/internal/model/user"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/api"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/page"
)

// plain text for unauthorized
func AuthJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := processJWT(w, r)
		if err != nil {
			api.ServeTextMessage(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// auth page for unauthorized
func AuthHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := processJWT(w, r)
		if err != nil {
			page.HandleAuth(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// check token validacy, generate new token if old one is expired and user still exists
func processJWT(w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.JWTClaimsByRequest(r)
	if err != auth.ErrExpToken {
		return err
	}

	// Get user id
	userIDfloat, ok := claims["userID"].(float64) // for some reason, in JWT it's stored as float64
	if !ok {
		return auth.ErrBadToken
	}

	// Check userID existance in database
	if !usermodel.IDExists(context.Background(), int(userIDfloat)) {
		return auth.ErrBadToken
	}

	// set new cookie
	cookie, err := auth.GenerateJWTCookie(int(userIDfloat), claims["username"].(string))
	if err != nil {
		return err
	}

	http.SetCookie(w, cookie)

	return nil
}
