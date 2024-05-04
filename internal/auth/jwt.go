package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTExp         = time.Hour * 24
	CookieLifeTime = time.Hour * 24 * 30
)

var (
	ErrBadToken = errors.New("invalid token")
	ErrExpToken = errors.New("expired token")
	secret      = []byte("")
)

func InitJWTSecret() {
	secret = []byte(env.GetEnvVal("JWT_SECRET"))
}

func GenerateJWTCookie(userID int, username string) (*http.Cookie, error) {
	token, err := createJWT(userID, username)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		MaxAge:   int(CookieLifeTime.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie, nil
}

func JWTClaimsByRequest(r *http.Request) (jwt.MapClaims, error) {
	// get token from cookie
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	// get token claims
	claims, err := jwtClaimsByString(cookie.Value)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func createJWT(userID int, username string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID":   userID,
			"username": username,
			"exp":      time.Now().Add(JWTExp).Unix(),
		},
	)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func jwtClaimsByString(tokenString string) (jwt.MapClaims, error) {
	token, err := jwtByString(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrBadToken
	}

	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if expTime.Before(time.Now()) {
		return nil, ErrExpToken
	}

	return claims, nil
}

func jwtByString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrBadToken
	}

	return token, nil
}
