package api

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/server/realtime/nodes/basenode"
	"github.com/1001bit/onlinecanvasgames/services/games/pkg/message"
)

func HandleRoomPost(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	gameTitle := r.PathValue("title")

	room, err := baseNode.ConnectNewRoom(r.Context(), gameTitle)
	if err != nil {
		ServeTextMessage(w, "Could not create a room!", http.StatusInternalServerError)
		return
	}

	resp := message.JSON{
		Type: "roomcreate",
		Body: room.GetID(),
	}

	ServeMessage(w, resp, http.StatusOK)
}
