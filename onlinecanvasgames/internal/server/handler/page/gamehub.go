package page

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/1001bit/OnlineCanvasGames/internal/gamemodel"
)

type GameHubData struct {
	Game *gamemodel.Game
}

func HandleGameHub(w http.ResponseWriter, r *http.Request) {
	data := GameHubData{}

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

	serveTemplate("gamehub.html", data, w, r)
}
