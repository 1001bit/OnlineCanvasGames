package gamelogic

import "github.com/1001bit/OnlineCanvasGames/internal/server/message"

type GameInfo struct {
	TicksPerSecond int `json:"tps"`
}

func NewGameInfoMessage(tps int) *message.JSON {
	return &message.JSON{
		Type: "gameinfo",
		Body: GameInfo{
			TicksPerSecond: tps,
		},
	}
}

func NewLevelMessage(l *Level) *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: l.GetPublicRects(),
	}
}
