package page

import (
	"net/http"
)

type GameHubData struct {
	GameTitle string
}

func HandleGameHub(w http.ResponseWriter, r *http.Request) {
	data := GameHubData{
		GameTitle: r.PathValue("title"),
	}

	serveTemplate(w, r, "gamehub.html", data)
}
