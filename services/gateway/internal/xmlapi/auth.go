package xmlapi

import (
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/components"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	if _, ok := claimscontext.GetUsername(r.Context()); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	components.Auth().Render(r.Context(), w)
}
