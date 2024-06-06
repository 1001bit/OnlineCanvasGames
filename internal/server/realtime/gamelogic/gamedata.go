package gamelogic

import "github.com/1001bit/OnlineCanvasGames/internal/server/message"

type GameData struct {
	TicksPerSecond int           `json:"tps"`
	Level          map[int]*Rect `json:"level"`
}

func NewGameDataMessage(tps int, l *Level) *message.JSON {
	return &message.JSON{
		Type: "gamedata",
		Body: GameData{
			TicksPerSecond: tps,
			Level:          l.GetPublicRects(),
		},
	}
}
