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

	createHubChan   chan *GameHub
	removeHubIDChan chan int
}

func NewGamesSSE() (*GamesSSE, error) {
	sse := &GamesSSE{
		SSELayer: MakeSSELayer(),
		hubs:     make(map[int]*GameHub),

		createHubChan:   make(chan *GameHub),
		removeHubIDChan: make(chan int),
	}

	return sse, nil
}

func (sse *GamesSSE) HandleSSE(w http.ResponseWriter, r *http.Request) {
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

	client := NewClient(w)
	hub.connect <- client

	client.writePump(r.Context().Done())
}

func (sse *GamesSSE) InitHubs() error {
	games, err := gamemodel.GetAll()
	if err != nil {
		return err
	}

	for _, game := range games {
		hub := NewGameHub()
		hub.id = game.ID
		go func() {
			sse.createHubChan <- hub
		}()
	}

	return nil
}

func (sse *GamesSSE) Run() {
	log.Println("<GameSSE Run>")

	for {
		select {
		// Client
		case client := <-sse.connect:
			sse.clients[client] = true
			log.Println("<GameSSE Client Connect>")

		case client := <-sse.disconnect:
			delete(sse.clients, client)
			log.Println("<GameSSE Client Disconnect>")

		// Hub
		case hub := <-sse.createHubChan:
			sse.hubs[hub.id] = hub
			hub.sse = sse
			go hub.Run()
			log.Println("<GameSSE Hub Create>")

		case hubID := <-sse.removeHubIDChan:
			delete(sse.hubs, hubID)
			log.Println("<GameSSE Hub Remove>")

		// Messages from server
		case message := <-sse.serverMessageChan:
			sse.handleServerMessage(message)
		}
	}
}

func (sse *GamesSSE) handleServerMessage(message string) {
	log.Println("<GameSSE Message>:", message)
}
