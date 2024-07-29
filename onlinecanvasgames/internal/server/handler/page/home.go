package page

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/claimscontext"
	"github.com/1001bit/OnlineCanvasGames/internal/gamemodel"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/nodes/basenode"
)

type HomeData struct {
	Username string
	Games    []gamemodel.Game
}

func HandleHome(w http.ResponseWriter, r *http.Request, baseNode *basenode.BaseNode) {
	data := HomeData{}

	_, username, _ := claimscontext.GetClaims(r.Context())
	data.Username = username

	// games list
	data.Games = baseNode.GetGamesJSON()

	serveTemplate("home.html", data, w, r)
}
