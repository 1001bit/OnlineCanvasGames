package page

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/gamesservice"
)

type GameHubData struct {
	GameTitle string
}

func HandleGameHub(w http.ResponseWriter, r *http.Request, gamesService *gamesservice.GamesService) {
	data := GameHubData{
		GameTitle: r.PathValue("title"),
	}

	serveTemplate(w, r, "gamehub.html", data)
}
