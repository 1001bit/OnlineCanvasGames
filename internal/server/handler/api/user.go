package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	usermodel "github.com/1001bit/OnlineCanvasGames/internal/model/user"
)

type UserPostRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

func HandleUserPost(w http.ResponseWriter, r *http.Request) {
	var request UserPostRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ServeJSONMessage(w, "Something went wrong. Please, try again", http.StatusBadRequest)
		return
	}

	// disallow empty fields
	if request.Password == "" || request.Username == "" {
		ServeJSONMessage(w, "Password or username field is empty", http.StatusBadRequest)
		return
	}

	// disallow username with special characters
	if request.Username != regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(request.Username, "") {
		ServeJSONMessage(w, "Username must not contain special characters", http.StatusBadRequest)
		return
	}

	// disallow short password
	if len(request.Password) < 8 {
		ServeJSONMessage(w, "Password should be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Login / register
	var user *usermodel.User
	switch request.Type {
	case "login":
		user, err = usermodel.GetByNameAndPassword(request.Username, request.Password)
	case "register":
		user, err = usermodel.Insert(request.Username, request.Password)
	}

	if err != nil {
		switch err {
		case usermodel.ErrNoSuchUser:
			ServeJSONMessage(w, "Incorrect username or password", http.StatusUnauthorized)
		case usermodel.ErrUserExists:
			ServeJSONMessage(w, fmt.Sprintf("%s already exists", request.Username), http.StatusUnauthorized)
		default:
			ServeJSONMessage(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("login/register err:", err)
		}
		return
	}

	// set token cookie
	token, err := auth.CreateJWT(user.ID, user.Name)
	if err != nil {
		ServeJSONMessage(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("jwt creation err:", err)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		MaxAge:   int(auth.JWTLifeTime.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	ServeJSONMessage(w, "Success!", http.StatusOK)
}
