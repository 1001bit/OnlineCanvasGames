package page

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/basenode"
)

type HomeData struct {
	Username string
	Games    []gamemodel.Game
}

func HandleHome(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	data := HomeData{}

	claims, err := auth.GetJwtClaimsFromContext(r.Context())
	if err == nil {
		data.Username = claims.Username
	}

	// games list
	data.Games = baseNode.GetGamesJSON()

	serveTemplate("home.html", data, w, r)
}
