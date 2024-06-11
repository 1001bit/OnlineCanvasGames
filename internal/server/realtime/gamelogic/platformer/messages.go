package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

type GameInfo struct {
	TPS    int `json:"tps"`
	Limit  int `json:"limit"`
	RectID int `json:"rectID"`
}

type LevelData struct {
	StaticRects    map[int]*physics.Rect          `json:"static"`
	KinematicRects map[int]*physics.KinematicRect `json:"kinematic"`
}

func (gl *PlatformerGL) NewLevelMessage() *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: LevelData{
			StaticRects:    gl.level.physEnv.GetStaticRects(),
			KinematicRects: gl.level.physEnv.GetKinematicRects(),
		},
	}
}

func (gl *PlatformerGL) NewDeleteMessage(rectID int) *message.JSON {
	return &message.JSON{
		Type: "delete",
		Body: rectID,
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
