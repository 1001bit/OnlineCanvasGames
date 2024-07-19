package accesstoken

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/basetoken"
)

type Claims struct {
	UserID   int
	Username string
}

// Get access token claims from request cookie
func GetClaims(r *http.Request) (Claims, error) {
	cookie, err := r.Cookie(Name)
	if err != nil {
		return Claims{}, err
	}

	// extract jwt claims from cookie string
	jwtClaims, err := basetoken.GetJwtClaims(cookie.Value)
	if err != nil {
		return Claims{}, err
	}

	// jwt claims -> claims
	claims := Claims{
		UserID:   int(jwtClaims["userID"].(float64)),
		Username: jwtClaims["username"].(string),
	}

	return claims, nil
}
