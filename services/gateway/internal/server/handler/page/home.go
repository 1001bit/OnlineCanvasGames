package page

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/auth/claimscontext"
	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/gamesservice"
)

type HomeData struct {
	Username string
	Titles   []string
}

func HandleHome(w http.ResponseWriter, r *http.Request, service *gamesservice.GamesService) {
	data := HomeData{}

	data.Username, _ = claimscontext.GetUsername(r.Context())

	data.Titles, _ = service.GetGames(r.Context())

	serveTemplate(w, r, "home.html", data)
}
