package gamesservice

import (
	"context"
	"net/http"

	"github.com/1001bit/onlinecanvasgames/services/gateway/pkg/auth/claimscontext"
	service "github.com/1001bit/onlinecanvasgames/services/gateway/pkg/client"
)

type Client struct {
	service *service.RestClient
}

func NewClient(host, port string) (*Client, error) {
	service, err := service.NewRestClient(host, port)
	if err != nil {
		return nil, err
	}
	return &Client{
		service: service,
	}, nil
}

func (s *Client) ProxyHandler() http.HandlerFunc {
	return s.service.Proxy()
}

func (s *Client) RoomProxyHandler() http.HandlerFunc {
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

func (s *Client) GetGames(ctx context.Context) ([]string, error) {
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
