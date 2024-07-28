package gamemodel

import (
	"context"
	"time"

	"github.com/1001bit/OnlineCanvasGames/internal/database"
)

const maxQueryTime = 5 * time.Second

type Game struct {
	ID    int
	Title string
}

func NewGame() *Game {
	return &Game{}
}

func GetAll(ctx context.Context) ([]Game, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx, "SELECT id, title FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game

	for rows.Next() {
		var game Game

		err := rows.Scan(&game.ID, &game.Title)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func GetByID(ctx context.Context, id int) (*Game, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	game := NewGame()
	game.ID = id

	err := database.DB.QueryRowContext(ctx, "SELECT title FROM games WHERE id = $1", id).Scan(&game.Title)
	if err != nil {
		return nil, err
	}

	return game, nil
}
