package api

import (
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime"
)

type RoomAPIResponse struct {
	RoomID int `json:"roomid"`
}

func HandleRoomPost(w http.ResponseWriter, r *http.Request, rt *realtime.Realtime) {
	resp := RoomAPIResponse{}

	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		ServeJSONMessage(w, "Bad game id", http.StatusBadRequest)
		return
	}

	room, err := rt.ConnectNewRoom(r.Context(), gameID)
	if err != nil {
		ServeJSONMessage(w, "Could not create a room!", http.StatusInternalServerError)
		return
	}
	resp.RoomID = room.GetID()

	ServeJSON(w, resp, http.StatusOK)
}
