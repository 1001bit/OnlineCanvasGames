package rt

import (
	"net/http"
	"strconv"

	basenode "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/base"
)

func HandleGameSSE(w http.ResponseWriter, r *http.Request, baseRT *basenode.BaseRT) {
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

	err = baseRT.ConnectToGame(r.Context(), w, gameID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
