package api

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/token"
	"github.com/1001bit/OnlineCanvasGames/internal/server/service"
)

// either register or login
func HandleUserPost(w http.ResponseWriter, r *http.Request, userService *service.UserService) {
	user, message := userService.PostUser(r.Context(), r.Body)

	if user == nil {
		ServeTextMessage(w, message, http.StatusBadRequest)
		return
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Name)
	if err != nil {
		ServeTextMessage(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	refreshToken, err := token.GenerateRefreshToken(user.ID, user.Name)
	if err != nil {
		ServeTextMessage(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	accessCookie := token.GenerateAccessTokenCookie(accessToken, false)
	refreshCookie := token.GenerateRefreshTokenCookie(refreshToken, false)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	ServeTextMessage(w, "Success!", http.StatusOK)
}
