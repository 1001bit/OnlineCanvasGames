package api

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/ws"
)

type RoomPostResponse struct {
	RoomID int `json:"roomid"`
}

func HandleRoomPost(w http.ResponseWriter, r *http.Request, ws *ws.GamesWS) {
	resp := RoomPostResponse{}

	room, err := ws.ConnectNewRoom(r.Context())
	if err != nil {
		ServeJSONMessage(w, "Could not create a room!", http.StatusInternalServerError)
		return
	}
	resp.RoomID = room.GetID()

	ServeJSON(w, resp, http.StatusOK)
}
