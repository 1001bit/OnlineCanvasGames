package basetoken

import "github.com/golang-jwt/jwt/v5"

func ClaimsToString(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
