package sse

import (
	"log"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

type GamesSSE struct {
	hubs map[int]*GameHub

	createHubChan   chan *GameHub
	removeHubIDChan chan int
}

func NewGamesSSE() (*GamesSSE, error) {
	sse := &GamesSSE{
		hubs: make(map[int]*GameHub),

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

	// TODO: Incorrect hubID error handling
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		return
	}

	hub, ok := sse.hubs[gameID]
	if !ok {
		return
	}

	client := NewClient(w)
	hub.connectChan <- client

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
		case hub := <-sse.createHubChan:
			sse.hubs[hub.id] = hub
			go hub.Run()

			log.Println("<GameSSE Hub Create>")

		case hubID := <-sse.removeHubIDChan:
			delete(sse.hubs, hubID)

			log.Println("<GameSSE Hub Remove>")
		}
	}
}
