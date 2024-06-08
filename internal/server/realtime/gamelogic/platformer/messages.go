package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/server/message"

func NewLevelMessage(l *Level) *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: l.GetPublicRects(),
	}
}
