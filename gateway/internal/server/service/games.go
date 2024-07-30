package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1001bit/ocg-gateway-service/internal/auth/claimscontext"
)

type Game struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type GamesService struct {
	service *Service
}

func NewGamesService(host, port string) (*GamesService, error) {
	service, err := NewService(host, port)
	if err != nil {
		return nil, err
	}
	return &GamesService{
		service: service,
	}, nil
}

func (s *GamesService) ProxyHandler() http.HandlerFunc {
	return s.service.proxy()
}

func (s *GamesService) RoomProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, username, err := claimscontext.GetClaims(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
		r.Header.Set("X-Username", username)

		s.service.proxy().ServeHTTP(w, r)
	}
}

func (s *GamesService) GetGameByID(ctx context.Context, id int) (*Game, error) {
	// forward post request from original body
	msg, err := s.service.request(ctx, "GET", fmt.Sprintf("/api/game/%d", id), nil)
	if err != nil {
		return nil, err
	}

	switch msg.Type {
	case "game":
		game := mapToGame(msg.Body.(map[string]any))
		if game == nil {
			return nil, ErrBadRequest
		}

		return game, nil
	default:
		return nil, ErrBadRequest
	}
}

func (s *GamesService) GetGames(ctx context.Context) ([]*Game, error) {
	// forward post request from original body
	msg, err := s.service.request(ctx, "GET", "/api/game", nil)
	if err != nil {
		return nil, err
	}

	switch msg.Type {
	case "games":
		// user type, ok
		body := msg.Body.([]any)

		games := make([]*Game, len(body))

		for i := range body {
			game := body[i].(map[string]any)
			if game == nil {
				return nil, ErrBadRequest
			}
			games[i] = mapToGame(game)
		}

		return games, nil
	default:
		return nil, ErrBadRequest
	}
}

func mapToGame(m map[string]any) *Game {
	game := &Game{}
	var ok bool

	id, ok := m["id"].(float64)
	if !ok {
		return nil
	}
	game.ID = int(id)

	game.Title, ok = m["title"].(string)
	if !ok {
		return nil
	}

	return game
}
