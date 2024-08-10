package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/token"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/userservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserLoginHandler(userService *userservice.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}
		json.NewDecoder(r.Body).Decode(&input)

		if message, ok := input.validate(); !ok {
			WriteTextMessage(w, message, http.StatusUnauthorized)
			return
		}

		user, err := userService.LoginUser(r.Context(), input.Username, input.Password)
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				HandleServerError(w)
				return
			}

			switch e.Code() {
			case codes.Unauthenticated:
				WriteTextMessage(w, "Invalid username or password", http.StatusUnauthorized)
			default:
				HandleServerError(w)
			}

			return
		}

		setTokenCookies(w, user.Username)
		WriteTextMessage(w, "Success!", http.StatusOK)
	}
}

func UserRegisterHandler(userService *userservice.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			HandleBadRequest(w)
			return
		}

		if message, ok := input.validate(); !ok {
			WriteTextMessage(w, message, http.StatusUnauthorized)
			return
		}

		user, err := userService.RegisterUser(r.Context(), input.Username, input.Password)
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				HandleServerError(w)
				return
			}

			switch e.Code() {
			case codes.AlreadyExists:
				WriteTextMessage(w, "User with such username already exists", http.StatusUnauthorized)
			default:
				HandleServerError(w)
			}

			return
		}

		setTokenCookies(w, user.Username)
		WriteTextMessage(w, "Success!", http.StatusOK)
	}
}

func (input *Input) validate() (string, bool) {
	// disallow empty fields
	if input.Password == "" || input.Username == "" {
		return "Password or username field is empty", false
	}

	// disallow username with special characters
	if input.Username != regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(input.Username, "") {
		return "Username must not contain special characters", false
	}

	// disallow short password
	if len(input.Password) < 8 {
		return "Password should be at least 8 characters long", false
	}

	return "", true
}

func setTokenCookies(w http.ResponseWriter, username string) {
	accessToken, err := token.GenerateAccessToken(username)
	if err != nil {
		WriteTextMessage(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	refreshToken, err := token.GenerateRefreshToken(username)
	if err != nil {
		WriteTextMessage(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	accessCookie := token.GenerateAccessTokenCookie(accessToken, false)
	refreshCookie := token.GenerateRefreshTokenCookie(refreshToken, false)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}
