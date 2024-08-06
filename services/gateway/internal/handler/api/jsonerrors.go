package api

import (
	"net/http"
)

func HandleUnauthorized(w http.ResponseWriter) {
	ServeTextMessage(w, "Unauthorized", http.StatusUnauthorized)
}

func HandleBadRequest(w http.ResponseWriter) {
	ServeTextMessage(w, "Bad request", http.StatusBadRequest)
}

func HandleServerError(w http.ResponseWriter) {
	ServeTextMessage(w, "Something went wrong!", http.StatusInternalServerError)
}
