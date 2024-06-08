package gamelogic

import "github.com/1001bit/OnlineCanvasGames/internal/server/message"

func NewGameInfoMessage(data any) *message.JSON {
	return &message.JSON{
		Type: "gameinfo",
		Body: data,
	}
}
