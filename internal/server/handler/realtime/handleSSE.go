package realtime

import (
	"net/http"
	"strconv"
)

// handle SSE endpoint
func (rt *Realtime) HandleGameSSE(w http.ResponseWriter, r *http.Request) {
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

	game, ok := rt.games.idMap[gameID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	client := NewGameRTClient(w)
	go client.Run(r.Context(), game)

	<-client.flow.Done()
}
