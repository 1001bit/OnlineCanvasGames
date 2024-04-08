package page

import (
	"database/sql"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
)

type GameRoomData struct {
	Game *gamemodel.Game
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

	serveTemplate("gameplay.html", data, w, r)
}
