package auth

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/accesstoken"
	"github.com/1001bit/OnlineCanvasGames/internal/auth/refreshtoken"
	usermodel "github.com/1001bit/OnlineCanvasGames/internal/model/user"
)

// Refresh access and refresh tokens by refresh token
func RefreshTokens(w http.ResponseWriter, r *http.Request) (accesstoken.Claims, error) {
	// Check user by refresh token
	refreshClaims, err := refreshtoken.ClaimsFromRequest(r)
	if err != nil {
		return accesstoken.Claims{}, err
	}
	// Check existance and get username
	user, err := usermodel.GetByID(r.Context(), refreshClaims.UserID)
	if err != nil {
		return accesstoken.Claims{}, err
	}

	// Generate cookies for tokens
	refreshCookie, err := refreshtoken.NewCookie(user.ID)
	if err != nil {
		return accesstoken.Claims{}, err
	}

	accessCookie, err := accesstoken.NewCookie(user.ID, user.Name)
	if err != nil {
		return accesstoken.Claims{}, err
	}

	// Set http cookies
	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)

	// Return access token claims
	claims := accesstoken.Claims{
		UserID:   user.ID,
		Username: user.Name,
	}
	return claims, nil
}
