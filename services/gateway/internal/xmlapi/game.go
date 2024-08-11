package xmlapi

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/components"
)

func HandleGameHub(w http.ResponseWriter, r *http.Request) {
	components.GameHub(r.PathValue("title")).Render(r.Context(), w)
}

func HandleGameRoom(w http.ResponseWriter, r *http.Request) {
	components.GameRoom(r.PathValue("roomid"), r.PathValue("title")).Render(r.Context(), w)
}
