package page

import (
	"net/http"
	"strconv"
)

type GameRoomData struct {
	RoomID int
}

func HandleGameRoom(w http.ResponseWriter, r *http.Request) {
	data := GameRoomData{}

	// roomID
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}

	data.RoomID = roomID

	serveTemplate("gameroom.html", data, w, r)
}
