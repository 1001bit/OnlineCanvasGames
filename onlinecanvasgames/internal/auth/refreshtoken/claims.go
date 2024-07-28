package refreshtoken

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/basetoken"
)

type Claims struct {
	UserID int
}

// extract claims from cookie in request
func ClaimsFromRequest(r *http.Request) (Claims, error) {
	cookie, err := r.Cookie(Name)
	if err != nil {
		return Claims{}, err
	}

	// Get jwt claims from string
	jwtClaims, err := basetoken.GetJwtClaims(cookie.Value)
	if err != nil {
		return Claims{}, err
	}

	// jwt claims -> claims
	claims := Claims{
		UserID: int(jwtClaims["userID"].(float64)),
	}

	return claims, nil
}
