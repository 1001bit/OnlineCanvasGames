package rt

import (
	"net/http"

	"github.com/1001bit/ocg-games-service/internal/server/realtime/nodes/basenode"
)

func HandleGameSSE(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get game from path
	gameTitle := r.PathValue("title")

	// connects client to sse (blocks the goroutine)
	err := baseNode.ConnectToGame(r.Context(), w, gameTitle)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
