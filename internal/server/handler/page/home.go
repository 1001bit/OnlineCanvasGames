package page

import (
	"fmt"
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

	claims, err := auth.JWTClaimsByRequest(r)
	if err == nil {
		data.Username = fmt.Sprint(claims["username"])
	}

	// games list
	data.Games = baseNode.GetGamesJSON()

	serveTemplate("home.html", data, w, r)
}
