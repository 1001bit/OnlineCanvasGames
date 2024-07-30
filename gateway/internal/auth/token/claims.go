package token

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
	UserID   int    `json:"userID"`
}
