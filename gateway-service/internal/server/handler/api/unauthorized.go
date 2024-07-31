package api

import (
	"net/http"
)

func HandleUnauthorized(w http.ResponseWriter, r *http.Request) {
	ServeTextMessage(w, "Unauthorized", http.StatusOK)
}
