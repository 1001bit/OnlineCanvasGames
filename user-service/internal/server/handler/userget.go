package handler

import (
	"database/sql"
	"net/http"

	"github.com/1001bit/ocg-user-service/internal/server/message"
	"github.com/1001bit/ocg-user-service/internal/usermodel"
)

func HandleUserGet(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("name")

	user, err := usermodel.GetByName(r.Context(), username)
	switch err {
	case nil:
		// continue
	case sql.ErrNoRows:
		ServeTextMessage(w, "No such user", http.StatusBadRequest)
		return
	default:
		ServeTextMessage(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	ServeMessage(w, message.JSON{
		Type: "user",
		Body: user,
	}, http.StatusOK)
}
