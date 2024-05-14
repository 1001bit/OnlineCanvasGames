package basetoken

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrBadToken = errors.New("invalid token")

func StringToClaims(tokenStr string) (jwt.MapClaims, error) {
	// Validate token
	token, err := stringToJWT(tokenStr)
	if err != nil {
		return nil, err
	}

	// get map claims
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrBadToken
	}

	return mapClaims, nil
}

func stringToJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
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