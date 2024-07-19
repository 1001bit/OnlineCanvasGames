package basetoken

import (
	"github.com/golang-jwt/jwt/v5"
)

// Extract jwt claims from string
func GetJwtClaims(tokenStr string) (jwt.MapClaims, error) {
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
