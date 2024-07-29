package page

import (
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/claimscontext"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	_, _, err := claimscontext.GetClaims(r.Context())

	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	serveTemplate("auth.html", nil, w, r)
}
