package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTExp         = time.Hour * 24 * 3  // 3d
	CookieLifeTime = time.Hour * 24 * 30 // 30d
)

func CookieFromClaims(claims Claims) (*http.Cookie, error) {
	tokenString, err := stringFromClaims(claims)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(CookieLifeTime.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie, nil
}

func stringFromClaims(claims Claims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID":   claims.UserID,
			"username": claims.Username,
			"exp":      time.Now().Add(JWTExp).Unix(),
		},
	)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
