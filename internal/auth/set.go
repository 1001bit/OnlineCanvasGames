package auth

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/accesstoken"
	"github.com/1001bit/OnlineCanvasGames/internal/auth/refreshtoken"
)

func SetTokens(w http.ResponseWriter, userID int, username string) error {
	accessCookie, err := accesstoken.ClaimsToCookie(userID, username)
	if err != nil {
		return err
	}

	refreshCookie, err := refreshtoken.ClaimsToCookie(userID)
	if err != nil {
		return err
	}

	// cookies
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	return nil
}
