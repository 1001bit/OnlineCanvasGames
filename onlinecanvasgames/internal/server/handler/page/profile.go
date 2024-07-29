package page

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/claimscontext"
	"github.com/1001bit/OnlineCanvasGames/internal/server/service"
)

type ProfileData struct {
	OwnerName string
	UserName  string
	Date      string
}

func HandleProfile(w http.ResponseWriter, r *http.Request, userService *service.UserService) {
	data := ProfileData{}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		HandleNotFound(w, r)
		return
	}

	user, err := userService.GetUserByID(r.Context(), id)
	switch err {
	case nil:
		// continue
	case context.DeadlineExceeded:
		HandleServerOverload(w, r)
		return
	default:
		HandleNotFound(w, r)
		return
	}

	data.OwnerName = user.Name
	data.Date, err = formatPostgresDate(user.Date)
	if err != nil {
		HandleServerError(w, r)
		return
	}

	_, username, _ := claimscontext.GetClaims(r.Context())
	data.UserName = username

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
