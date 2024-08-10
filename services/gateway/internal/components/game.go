package components

import "net/http"

func HandleGameHub(w http.ResponseWriter, r *http.Request) {
	GameHub(r.PathValue("title")).Render(r.Context(), w)
}

func HandleGameRoom(w http.ResponseWriter, r *http.Request) {
	GameRoom(r.PathValue("roomid"), r.PathValue("title")).Render(r.Context(), w)
}
