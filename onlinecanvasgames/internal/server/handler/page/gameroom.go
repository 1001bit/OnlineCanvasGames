package page

import (
	"context"
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/server/service"
)

type GameRoomData struct {
	RoomID int
	Game   *service.Game
}

func HandleGameRoom(w http.ResponseWriter, r *http.Request, gamesService *service.GamesService) {
	data := GameRoomData{}

	// roomID
	roomID, err := strconv.Atoi(r.PathValue("roomid"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}
	data.RoomID = roomID

	// Game
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}

	data.Game, err = gamesService.GetGameByID(r.Context(), gameID)

	switch err {
	case nil:
		// continue
	case context.DeadlineExceeded:
		HandleServerOverload(w, r)
		return
	default:
		HandleNotFound(w, r)
		return
	}

	serveTemplate("gameroom.html", data, w, r)
}
