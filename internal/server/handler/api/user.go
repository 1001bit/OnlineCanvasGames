package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	usermodel "github.com/1001bit/OnlineCanvasGames/internal/model/user"
)

var ErrBadInput = errors.New("bad auth input")

type UserPostRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

func HandleUserPost(w http.ResponseWriter, r *http.Request) {
	var request UserPostRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ServeTextMessage(w, "Something went wrong. Please, try again", http.StatusBadRequest)
		return
	}

	if text, err := validateAuthInput(request.Username, request.Password); err != nil {
		ServeTextMessage(w, text, http.StatusBadRequest)
		return
	}

	// Login / register
	var user *usermodel.User
	switch request.Type {
	case "login":
		user, err = usermodel.GetByNameAndPassword(r.Context(), request.Username, request.Password)
	case "register":
		user, err = usermodel.Insert(r.Context(), request.Username, request.Password)
	}

	if err != nil {
		switch err {
		case usermodel.ErrNoSuchUser:
			ServeTextMessage(w, "Incorrect username or password", http.StatusUnauthorized)
		case usermodel.ErrUserExists:
			ServeTextMessage(w, fmt.Sprintf("%s already exists", request.Username), http.StatusUnauthorized)
		case context.DeadlineExceeded:
			ServeTextMessage(w, "Deadline exceeded", http.StatusInternalServerError)
			log.Println("Auth deadline exceeded", err)
		default:
			ServeTextMessage(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("login/register err:", err)
		}
		return
	}

	err = auth.SetTokens(w, user.ID, user.Name)
	if err != nil {
		ServeTextMessage(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("jwt creation err:", err)
		return
	}

	ServeTextMessage(w, "Success!", http.StatusOK)
}

func validateAuthInput(username, password string) (string, error) {
	// disallow empty fields
	if password == "" || username == "" {
		return "Password or username field is empty", ErrBadInput
	}

	// disallow username with special characters
	if username != regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(username, "") {
		return "Username must not contain special characters", ErrBadInput
	}

	// disallow short password
	if len(password) < 8 {
		return "Password should be at least 8 characters long", ErrBadInput
	}

	return "", nil
}
