package page

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	_, ok := claimscontext.GetUsername(r.Context())

	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	serveTemplate(w, r, "auth.html", nil)
}
