package page

import (
	"net/http"
	"strconv"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/gamesservice"
)

type GameRoomData struct {
	RoomID    int
	GameTitle string
}

func HandleGameRoom(w http.ResponseWriter, r *http.Request, gamesService *gamesservice.GamesService) {
	data := GameRoomData{
		GameTitle: r.PathValue("title"),
	}

	// roomID
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}
	data.RoomID = roomID

	serveTemplate(w, r, "gameroom.html", data)
}
