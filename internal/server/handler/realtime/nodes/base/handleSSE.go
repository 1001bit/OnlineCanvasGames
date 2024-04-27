package basenode

import (
	"net/http"
	"strconv"

	gamenode "github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime/nodes/game"
)

// handle SSE endpoint
func (baseRT *BaseRT) HandleGameSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get game from path
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gameRT, ok := baseRT.games.IDMap[gameID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	client := gamenode.NewGameRTClient(w)

	// RUN GameRTClient
	go func() {
		gameRT.Clients.ConnectChild(client)
		client.Run(r.Context())
		gameRT.Clients.DisconnectChild(client)
	}()

	<-client.Flow.Done()
}
