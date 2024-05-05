package auth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var ErrBadToken = errors.New("invalid token")

type Claims struct {
	UserID   int
	Username string
}

func ClaimsFromRequest(r *http.Request) (Claims, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return Claims{}, err
	}

	return claimsFromString(cookie.Value)
}

func claimsFromString(tokenString string) (Claims, error) {
	token, err := stringToJWT(tokenString)
	if err != nil {
		return Claims{}, err
	}

	return claimsFromJWT(token)
}

func claimsFromJWT(token *jwt.Token) (Claims, error) {
	// get map claims
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, ErrBadToken
	}

	claims := Claims{
		UserID:   int(mapClaims["userID"].(float64)),
		Username: mapClaims["username"].(string),
	}

	return claims, nil
}

func stringToJWT(tokenString string) (*jwt.Token, error) {
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
