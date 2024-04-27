package api

import (
	"net/http"
	"strconv"

	rtnode "github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime/node"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

func HandleRoomPost(w http.ResponseWriter, r *http.Request, baseRT *rtnode.BaseRT) {
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		ServeTextMessage(w, "Bad game id", http.StatusBadRequest)
		return
	}

	room, err := baseRT.ConnectNewRoom(r.Context(), gameID)
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
