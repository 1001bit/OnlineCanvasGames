package xmlapi

import (
	"net/http"
	"regexp"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/components"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/token"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/userservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoginHandler(serviceClient *userservice.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if message, ok := validateInput(username, password); !ok {
			renderMessage(w, r, message, http.StatusBadRequest)
			return
		}

		user, err := serviceClient.LoginUser(r.Context(), username, password)
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				renderMessage(w, r, "Something went wrong!", http.StatusInternalServerError)
				return
			}

			switch e.Code() {
			case codes.Unauthenticated:
				renderMessage(w, r, "Invalid username or password", http.StatusUnauthorized)
			default:
				renderMessage(w, r, "Something went wrong!", http.StatusInternalServerError)
			}

			return
		}

		ok := setTokenCookies(w, user.Username)
		if !ok {
			renderMessage(w, r, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		w.Header().Add("HX-Redirect", "/")
		renderMessage(w, r, "Success!", http.StatusOK)
	}
}

func RegisterHandler(serviceClient *userservice.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if message, ok := validateInput(username, password); !ok {
			renderMessage(w, r, message, http.StatusBadRequest)
			return
		}

		user, err := serviceClient.RegisterUser(r.Context(), username, password)
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				renderMessage(w, r, "Something went wrong!", http.StatusInternalServerError)
				return
			}

			switch e.Code() {
			case codes.AlreadyExists:
				renderMessage(w, r, "User with such username already exists", http.StatusUnauthorized)
			default:
				renderMessage(w, r, "Something went wrong!", http.StatusInternalServerError)
			}

			return
		}

		ok := setTokenCookies(w, user.Username)
		if !ok {
			renderMessage(w, r, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		w.Header().Add("HX-Redirect", "/")
		renderMessage(w, r, "Success!", http.StatusOK)
	}
}

func renderMessage(w http.ResponseWriter, r *http.Request, text string, code int) {
	w.WriteHeader(code)
	components.AuthInfo(text).Render(r.Context(), w)
}

func validateInput(username, password string) (string, bool) {
	// disallow empty fields
	if password == "" || username == "" {
		return "Password or username field is empty", false
	}

	// disallow username with special characters
	if username != regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(username, "") {
		return "Username must not contain special characters", false
	}

	// disallow short password
	if len(password) < 8 {
		return "Password should be at least 8 characters long", false
	}

	return "", true
}

func setTokenCookies(w http.ResponseWriter, username string) bool {
	accessToken, err := token.GenerateAccessToken(username)
	if err != nil {
		return false
	}

	refreshToken, err := token.GenerateRefreshToken(username)
	if err != nil {
		return false
	}

	accessCookie := token.GenerateAccessTokenCookie(accessToken, false)
	refreshCookie := token.GenerateRefreshTokenCookie(refreshToken, false)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	return true
}
