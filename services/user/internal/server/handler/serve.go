package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/user/internal/server/message"
)

func ServeMessage(w http.ResponseWriter, msg message.JSON, status int) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println("err on response:", err)
		ServeTextMessage(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(msgByte)
}

func ServeTextMessage(w http.ResponseWriter, text string, status int) {
	msg := message.JSON{
		Type: "message",
		Body: text,
	}
	ServeMessage(w, msg, status)
}
