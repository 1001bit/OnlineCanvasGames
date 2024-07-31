package page

import (
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/server/service"
)

type GameHubData struct {
	GameTitle string
}

func HandleGameHub(w http.ResponseWriter, r *http.Request, gamesService *service.GamesService) {
	data := GameHubData{
		GameTitle: r.PathValue("title"),
	}

	serveTemplate("gamehub.html", data, w, r)
}
