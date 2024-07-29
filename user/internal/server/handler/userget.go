package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/neinBit/ocg-user-service/internal/server/message"
	"github.com/neinBit/ocg-user-service/internal/usermodel"
)

func HandleUserGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ServeTextMessage(w, "No such user", http.StatusBadRequest)
		return
	}

	user, err := usermodel.GetByID(r.Context(), id)
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
