package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/realtime"
)

type RoomPostResponse struct {
	RoomID int `json:"roomid"`
}

type RoomPostRequest struct {
	GameID int `json:"gameid"`
}

func HandleRoomPost(w http.ResponseWriter, r *http.Request, games *realtime.Realtime) {
	resp := RoomPostResponse{}
	req := RoomPostRequest{}
	json.NewDecoder(r.Body).Decode(&req)

	room, err := games.ConnectNewRoom(r.Context(), req.GameID)
	if err != nil {
		log.Println(req.GameID)
		ServeJSONMessage(w, "Could not create a room!", http.StatusInternalServerError)
		return
	}
	resp.RoomID = room.GetID()

	ServeJSON(w, resp, http.StatusOK)
}
