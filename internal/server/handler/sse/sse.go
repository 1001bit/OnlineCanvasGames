package sse

import (
	"log"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

type GamesSSE struct {
	SSELayer
	hubs map[int]*GameHub
}

func NewGamesSSE() (*GamesSSE, error) {
	games, err := gamemodel.GetAll()
	if err != nil {
		return nil, err
	}

	sse := &GamesSSE{
		SSELayer: MakeSSELayer(),
		hubs:     make(map[int]*GameHub),
	}

	for _, game := range games {
		sse.hubs[game.ID] = NewGameHub(sse)
		go sse.hubs[game.ID].Run()
	}

	return sse, nil
}

func (sse *GamesSSE) HandleEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return
	}
	hub, ok := sse.hubs[id]
	if !ok {
		return
	}

	client := NewClient(w, hub)
	client.hub.connect <- client
	client.writePump(r.Context().Done())
}

func (sse *GamesSSE) Run() {
	for {
		select {
		case client := <-sse.connect:
			sse.clients[client] = true
			log.Println("<GameSSE Connect>")

		case client := <-sse.disconnect:
			delete(sse.clients, client)
			log.Println("<GameSSE Disconnect>")

		case message := <-sse.messageChan:
			sse.handleMessage(message)
		}
	}
}

func (sse *GamesSSE) handleMessage(message string) {
	log.Println("<GameSSE Message>:", message)
}
