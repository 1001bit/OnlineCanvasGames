package page

import (
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/auth/claimscontext"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	_, ok := claimscontext.GetUsername(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	serveTemplate("auth.html", nil, w, r)
}
