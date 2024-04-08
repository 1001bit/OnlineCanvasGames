package page

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/handler/ws"
)

type GameRoomData struct {
	Game   *gamemodel.Game
	RoomID int
}

func HandleGamePlay(ws *ws.GamesWS, w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, fmt.Sprintf("/game/%d/play?room=%d", data.Game.ID, ws.PickRandomRoomID()), http.StatusSeeOther)
	case "-1":
		// TODO: No room exists
		return
	default:
		data.RoomID, err = strconv.Atoi(r.URL.Query().Get("room"))
		if err != nil {
			// TODO: Room not found
			return
		}
	}

	serveTemplate("gameplay.html", data, w, r)
}
