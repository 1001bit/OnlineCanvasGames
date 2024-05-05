package page

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	// only unauthorized may see this page
	_, ok := r.Context().Value(auth.ClaimsKey).(auth.Claims)

	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	serveTemplate("auth.html", nil, w, r)
}
