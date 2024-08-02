package token

import (
	"net/http"
	"time"

	"github.com/1001bit/overenv"
	"github.com/golang-jwt/jwt/v5"
)

var (
	refreshTokenSecret   = []byte(overenv.Get("TOKEN_SECRET"))
	refreshTokenDuration = time.Hour * 24 * 7
)

func GenerateRefreshToken(username string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
		},

		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

func GenerateRefreshTokenCookie(tokenString string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(refreshTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, refreshTokenSecret)
}
