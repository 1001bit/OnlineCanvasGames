package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/server/message"

type GameInfo struct {
	TPS    int `json:"tps"`
	Limit  int `json:"limit"`
	RectID int `json:"rectID"`
}

func (gl *PlatformerGL) NewFullLevelMessage() *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: gl.level.GetPublicRects(),
	}
}

func (gl *PlatformerGL) NewGameInfoMessage(playerRectID int) *message.JSON {
	return &message.JSON{
		Type: "gameinfo",
		Body: GameInfo{
			TPS:    gl.tps,
			Limit:  gl.maxPlayers,
			RectID: playerRectID,
		},
	}
}
