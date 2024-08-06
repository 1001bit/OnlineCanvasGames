package page

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/gamesservice"
)

type HomeData struct {
	Username string
	Titles   []string
}

func HandleHome(w http.ResponseWriter, r *http.Request, service *gamesservice.Client) {
	data := HomeData{}

	data.Username, _ = claimscontext.GetUsername(r.Context())

	data.Titles, _ = service.GetGames(r.Context())

	serveTemplate(w, r, "home.html", data)
}
