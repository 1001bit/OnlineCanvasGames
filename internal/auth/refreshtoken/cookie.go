package refreshtoken

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/basetoken"
	"github.com/golang-jwt/jwt/v5"
)

// Generate refresh token cookie by userID
func NewCookie(userID int) (*http.Cookie, error) {
	claims := Claims{
		UserID: userID,
	}

	tokenStr, err := claimsToString(claims)
	if err != nil {
		return nil, err
	}

	return basetoken.NewCookie(tokenStr, Name, ExpTime), nil
}

// Claims -> jwt Claims -> string
func claimsToString(claims Claims) (string, error) {
	mapClaims := jwt.MapClaims{
		"userID": claims.UserID,
		"exp":    ExpTime,
	}

	tokenStr, err := basetoken.JwtClaimsToString(mapClaims)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
