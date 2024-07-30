package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1001bit/OnlineCanvasGames/internal/auth/claimscontext"
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

func (s *GamesService) HandleRoomPost() http.HandlerFunc {
	return s.service.proxy()
}

func (s *GamesService) HandleGameHubSSE() http.HandlerFunc {
	return s.service.proxy()
}

func (s *GamesService) HandleRoomWS() http.HandlerFunc {
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
		// user type, ok
		body := msg.Body.(map[string]any)
		var ok bool

		game := &Game{
			ID: id,
		}

		game.Title, ok = body["title"].(string)
		if !ok {
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
			game := &Game{}

			// get i element of received body
			gameMap, ok := body[i].(map[string]any)
			if !ok {
				return nil, ErrBadRequest
			}

			id, ok := gameMap["id"].(float64)
			if !ok {
				return nil, ErrBadRequest
			}
			game.ID = int(id)

			game.Title, ok = gameMap["title"].(string)
			if !ok {
				return nil, ErrBadRequest
			}

			games[i] = game
		}

		return games, nil
	default:
		return nil, ErrBadRequest
	}
}