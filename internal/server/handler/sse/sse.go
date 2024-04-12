package sse

import (
	"log"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

type GamesSSE struct {
	hubs             map[int]*GameHub
	connectHubChan   chan *GameHub
	connectHubIDChan chan int
}

func NewGamesSSE() (*GamesSSE, error) {
	sse := &GamesSSE{
		hubs: make(map[int]*GameHub),

		connectHubChan:   make(chan *GameHub),
		connectHubIDChan: make(chan int),
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
	hub.connectClientChan <- client

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
			log.Println("<GameSSE Hub Connect>")

		case hubID := <-sse.connectHubIDChan:
			sse.disconnectHubByID(hubID)
			log.Println("<GameSSE Hub Disconnect>")
		}
	}
}

func (sse *GamesSSE) connectHub(hub *GameHub) {
	sse.hubs[hub.id] = hub
	go hub.Run()
}

func (sse *GamesSSE) disconnectHubByID(id int) {
	delete(sse.hubs, id)
}
