package auth

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/accesstoken"
	"github.com/1001bit/OnlineCanvasGames/internal/auth/refreshtoken"
	usermodel "github.com/1001bit/OnlineCanvasGames/internal/model/user"
)

func RefreshTokens(w http.ResponseWriter, r *http.Request) (accesstoken.Claims, error) {
	// userID
	refreshClaims, err := refreshtoken.ClaimsFromRequest(r)
	if err != nil {
		return accesstoken.Claims{}, err
	}
	userID := refreshClaims.UserID

	user, err := usermodel.GetByID(r.Context(), userID)
	if err != nil {
		return accesstoken.Claims{}, err
	}

	// cookies
	refreshCookie, err := refreshtoken.ClaimsToCookie(user.ID)
	if err != nil {
		return accesstoken.Claims{}, err
	}

	accessCookie, err := accesstoken.ClaimsToCookie(user.ID, user.Name)
	if err != nil {
		return accesstoken.Claims{}, err
	}

	// set cookies
	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)

	claims := accesstoken.Claims{
		UserID:   user.ID,
		Username: user.Name,
	}

	return claims, nil
}
