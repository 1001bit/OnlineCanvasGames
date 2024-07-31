package page

import (
	"net/http"
	"strconv"

	"github.com/1001bit/ocg-gateway-service/internal/server/service"
)

type GameRoomData struct {
	RoomID    int
	GameTitle string
}

func HandleGameRoom(w http.ResponseWriter, r *http.Request, gamesService *service.GamesService) {
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
