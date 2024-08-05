package gamesservice

import (
	"context"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/internal/server/service"
	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
)

type GamesService struct {
	service *service.Rest
}

func New(host, port string) (*GamesService, error) {
	service, err := service.NewRestService(host, port)
	if err != nil {
		return nil, err
	}
	return &GamesService{
		service: service,
	}, nil
}

func (s *GamesService) ProxyHandler() http.HandlerFunc {
	return s.service.Proxy()
}

func (s *GamesService) RoomProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := claimscontext.GetUsername(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-Username", username)

		s.ProxyHandler().ServeHTTP(w, r)
	}
}

func (s *GamesService) GetGames(ctx context.Context) ([]string, error) {
	// forward post request from original body
	msg, err := s.service.Request(ctx, "GET", "/api/game", nil)
	if err != nil {
		return nil, err
	}

	switch msg.Type {
	case "games":
		// user type, ok
		body := msg.Body.([]any)

		games := make([]string, len(body))

		for i := range body {
			game := body[i].(map[string]any)
			if game == nil {
				return nil, service.ErrBadRequest
			}
			games[i], _ = game["title"].(string)
		}

		return games, nil
	default:
		return nil, service.ErrBadRequest
	}
}
