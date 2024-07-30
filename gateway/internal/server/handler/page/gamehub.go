package page

import (
	"context"
	"net/http"
	"strconv"

	"github.com/1001bit/ocg-gateway-service/internal/server/service"
)

type GameHubData struct {
	Game *service.Game
}

func HandleGameHub(w http.ResponseWriter, r *http.Request, gamesService *service.GamesService) {
	data := GameHubData{}

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

	serveTemplate("gamehub.html", data, w, r)
}
