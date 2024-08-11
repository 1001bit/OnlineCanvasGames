package gamestore

import (
	"encoding/json"
	"os"
)

type Game struct {
	Title string `json:"title"`
}

func GetAllGames(jsonFileName string) ([]Game, error) {
	data, err := os.ReadFile(jsonFileName)
	if err != nil {
		return nil, err
	}

	var games []Game
	err = json.Unmarshal(data, &games)
	if err != nil {
		return nil, err
	}

	return games, nil
}
