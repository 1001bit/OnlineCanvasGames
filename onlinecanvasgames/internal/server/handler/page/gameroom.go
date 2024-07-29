package page

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/gamemodel"
)

type GameRoomData struct {
	RoomID int
	Game   *gamemodel.Game
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

	// Game
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}

	data.Game, err = gamemodel.GetByID(r.Context(), gameID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			HandleNotFound(w, r)
		case context.DeadlineExceeded:
			HandleServerOverload(w, r)
		default:
			HandleServerError(w, r)
		}
		return
	}

	serveTemplate("gameroom.html", data, w, r)
}
