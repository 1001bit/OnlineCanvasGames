package accesstoken

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/basetoken"
)

type Claims struct {
	UserID   int
	Username string
}

func ClaimsFromRequest(r *http.Request) (Claims, error) {
	cookie, err := r.Cookie(Name)
	if err != nil {
		return Claims{}, err
	}

	mapClaims, err := basetoken.StringToClaims(cookie.Value)
	if err != nil {
		return Claims{}, err
	}

	claims := Claims{
		UserID:   int(mapClaims["userID"].(float64)),
		Username: mapClaims["username"].(string),
	}

	return claims, nil
}
