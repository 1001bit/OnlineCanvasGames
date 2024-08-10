package components

import (
	"net/http"
	"time"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client/userservice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileData struct {
	OwnerName string
	UserName  string
	Date      string
}

func ProfileHandler(userService *userservice.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		selfUsername, _ := claimscontext.GetUsername(r.Context())
		username := r.PathValue("name")

		user, err := userService.GetUserByName(r.Context(), username)
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				ErrorInternal().Render(r.Context(), w)
				return
			}

			switch e.Code() {
			case codes.NotFound:
				ErrorNotFound().Render(r.Context(), w)
			default:
				ErrorInternal().Render(r.Context(), w)
			}

			return
		}

		username = user.Username
		date, _ := formatPostgresDate(user.Date)

		Profile(selfUsername, username, date).Render(r.Context(), w)
	}
}

// 2006-01-02T15:04:05Z -> 2 January 2006
func formatPostgresDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("2 January 2006"), nil
}
