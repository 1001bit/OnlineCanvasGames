package api

import (
	"net/http"
)

func HandleUnauthorized(w http.ResponseWriter) {
	WriteTextMessage(w, "Unauthorized", http.StatusUnauthorized)
}

func HandleBadRequest(w http.ResponseWriter) {
	WriteTextMessage(w, "Bad request", http.StatusBadRequest)
}

func HandleServerError(w http.ResponseWriter) {
	WriteTextMessage(w, "Something went wrong!", http.StatusInternalServerError)
}
