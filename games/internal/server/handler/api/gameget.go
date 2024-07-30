package api

import (
	"net/http"
	"strconv"

	"github.com/neinBit/ocg-games-service/internal/server/message"
	"github.com/neinBit/ocg-games-service/internal/server/realtime/nodes/basenode"
)

func HandleGameGet(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		ServeTextMessage(w, "Bad game id", http.StatusBadRequest)
		return
	}

	game, ok := baseNode.GetGameByID(gameID)
	if !ok {
		ServeTextMessage(w, "Not found", http.StatusNotFound)
		return
	}

	ServeMessage(w, message.JSON{
		Type: "game",
		Body: game.GetGame(),
	}, http.StatusOK)
}

func HandleGamesGet(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	games := baseNode.GetGamesJSON()

	ServeMessage(w, message.JSON{
		Type: "games",
		Body: games,
	}, http.StatusOK)
}
