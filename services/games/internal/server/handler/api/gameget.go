package api

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/nodes/basenode"
	"github.com/1001bit/onlinecanvasgames/services/games/pkg/message"
)

func HandleGamesGet(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	games := baseNode.GetGamesJSON()

	ServeMessage(w, message.JSON{
		Type: "games",
		Body: games,
	}, http.StatusOK)
}
