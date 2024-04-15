package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func ServeJSON(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)
	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Println("err on response:", err)
		ServeJSONMessage(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.Write(dataByte)
}

func ServeJSONMessage(w http.ResponseWriter, message string, status int) {
	data := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	ServeJSON(w, data, status)
}
