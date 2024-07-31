package page

import (
	"context"
	"net/http"
	"time"

	"github.com/1001bit/ocg-gateway-service/internal/auth/claimscontext"
	"github.com/1001bit/ocg-gateway-service/internal/server/service"
)

type ProfileData struct {
	OwnerName string
	UserName  string
	Date      string
}

func HandleProfile(w http.ResponseWriter, r *http.Request, userService *service.UserService) {
	data := ProfileData{}
	data.UserName, _ = claimscontext.GetUsername(r.Context())

	name := r.PathValue("name")

	user, err := userService.GetUserByName(r.Context(), name)
	switch err {
	case nil:
		// continue
	case context.DeadlineExceeded:
		HandleServerOverload(w, r)
		return
	default:
		if data.UserName == name {
			HandleLogout(w, r)
			return
		}

		HandleNotFound(w, r)
		return
	}

	data.OwnerName = user.Name
	data.Date, err = formatPostgresDate(user.Date)
	if err != nil {
		HandleServerError(w, r)
		return
	}

	serveTemplate("profile.html", data, w, r)
}

// 2006-01-02T15:04:05Z -> 2 January 2006
func formatPostgresDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("2 January 2006"), nil
}
