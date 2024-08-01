package token

import (
	"net/http"
	"time"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/env"
	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenSecret   = []byte(env.GetEnvVal("TOKEN_SECRET"))
	accessTokenDuration = time.Minute * 15
)

func GenerateAccessToken(username string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
		},

		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

func GenerateAccessTokenCookie(tokenString string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     "access",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(accessTokenDuration.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, accessTokenSecret)
}
