package sse

import (
	"context"
	"log"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

type GamesSSE struct {
	hubs              map[int]*GameHub
	connectHubChan    chan *GameHub
	disconnectHubChan chan *GameHub
}

func NewGamesSSE() (*GamesSSE, error) {
	sse := &GamesSSE{
		hubs: make(map[int]*GameHub),

		connectHubChan:    make(chan *GameHub),
		disconnectHubChan: make(chan *GameHub),
	}

	return sse, nil
}

func (sse *GamesSSE) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hub, ok := sse.hubs[gameID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	client := NewClient(w)
	hub.connectClientChan <- client

	client.writePump(r.Context().Done())
}

func (sse *GamesSSE) InitHubs() error {
	games, err := gamemodel.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, game := range games {
		hub := NewGameHub()
		hub.gameID = game.ID
		go func() {
			sse.connectHubChan <- hub
		}()
	}

	return nil
}

func (sse *GamesSSE) Run() {
	log.Println("<GameSSE Run>")
	defer log.Println("<GameSSE Run End>")

	for {
		select {
		case hub := <-sse.connectHubChan:
			sse.connectHub(hub)
			log.Println("<GameSSE +Hub>:", len(sse.hubs))

		case hub := <-sse.disconnectHubChan:
			sse.disconnectHub(hub)
			log.Println("<GameSSE -Hub>:", len(sse.hubs))
		}
	}
}

func (sse *GamesSSE) connectHub(hub *GameHub) {
	sse.hubs[hub.gameID] = hub
	go hub.Run()
}

func (sse *GamesSSE) disconnectHub(hub *GameHub) {
	delete(sse.hubs, hub.gameID)
}
