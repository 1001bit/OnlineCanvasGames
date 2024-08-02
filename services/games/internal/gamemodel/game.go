package gamemodel

import (
	"context"
	"database/sql"
)

type Game struct {
	Title string `json:"title"`
}

type GameStore struct {
	db *sql.DB
}

func NewGameStore(db *sql.DB) *GameStore {
	return &GameStore{
		db: db,
	}
}

func (gs *GameStore) GetAllGames(ctx context.Context) ([]Game, error) {
	rows, err := gs.db.QueryContext(ctx, "SELECT title FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]Game, 0)

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
