package page

import (
	"net/http"
	"time"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service/userservice"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileData struct {
	OwnerName string
	UserName  string
	Date      string
}

func HandleProfile(w http.ResponseWriter, r *http.Request, userService *userservice.UserService) {
	data := ProfileData{}
	data.UserName, _ = claimscontext.GetUsername(r.Context())

	name := r.PathValue("name")

	user, err := userService.GetUserByName(r.Context(), name)
	if err != nil {
		e, ok := status.FromError(err)
		if !ok {
			HandleServerError(w, r)
			return
		}

		switch e.Code() {
		case codes.NotFound:
			HandleNotFound(w, r)
		default:
			HandleServerError(w, r)
		}

		return
	}

	data.OwnerName = user.Username
	data.Date, err = formatPostgresDate(user.Date)
	if err != nil {
		HandleServerError(w, r)
		return
	}

	serveTemplate(w, r, "profile.html", data)
}

// 2006-01-02T15:04:05Z -> 2 January 2006
func formatPostgresDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("2 January 2006"), nil
}
