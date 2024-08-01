package gamemodel

import (
	"context"
	"time"

	"github.com/1001bit/onlinecanvasgames/services/games/internal/database"
)

const maxQueryTime = 5 * time.Second

type Game struct {
	Title string `json:"title"`
}

func NewGame() *Game {
	return &Game{}
}

func GetAll(ctx context.Context) ([]Game, error) {
	ctx, cancel := context.WithTimeout(ctx, maxQueryTime)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx, "SELECT title FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game

	for rows.Next() {
		var game Game

		err := rows.Scan(&game.Title)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}
