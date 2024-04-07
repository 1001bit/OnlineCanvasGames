package page

import (
	"database/sql"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

type GameRoomData struct {
	Game   *gamemodel.Game
	RoomID int
}

func HandleGamePlay(w http.ResponseWriter, r *http.Request) {
	data := GameRoomData{}

	// GameID
	gameID, err := strconv.Atoi(r.PathValue("gameid"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}

	data.Game, err = gamemodel.GetByID(gameID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			HandleNotFound(w, r)
		default:
			HandleServerError(w, r)
		}
		return
	}

	// RoomID
	switch r.URL.Query().Get("room") {
	case "":
		// TODO: redirect to random room
		data.RoomID = 0
	default:
		data.RoomID, err = strconv.Atoi(r.URL.Query().Get("room"))
		if err != nil {
			// TODO: Room not found
			return
		}
	}

	serveTemplate("gameplay.html", data, w, r)
}
