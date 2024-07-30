package page

import (
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/auth/claimscontext"
	"github.com/1001bit/ocg-gateway-service/internal/server/service"
)

type HomeData struct {
	Username string
	Games    []*service.Game
}

func HandleHome(w http.ResponseWriter, r *http.Request, service *service.GamesService) {
	data := HomeData{}

	_, username, _ := claimscontext.GetClaims(r.Context())
	data.Username = username

	data.Games, _ = service.GetGames(r.Context())

	serveTemplate("home.html", data, w, r)
}
