package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/message"
)

func WriteJsonMessage(w http.ResponseWriter, msg message.JSON, status int) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println("err on response:", err)
		WriteTextMessage(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(msgByte)
}

func WriteTextMessage(w http.ResponseWriter, text string, status int) {
	msg := message.JSON{
		Type: "message",
		Body: text,
	}
	WriteJsonMessage(w, msg, status)
}
