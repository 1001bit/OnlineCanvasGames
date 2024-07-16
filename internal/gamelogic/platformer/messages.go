package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

type GameInfo struct {
	PlayerRectID int       `json:"rectID"`
	Constants    Constants `json:"constants"`
	TPS          int       `json:"tps"`
}

type LevelData struct {
	StaticRects    map[int]*physics.Rect          `json:"static"`
	KinematicRects map[int]*physics.KinematicRect `json:"kinematic"`
}

type CreateData struct {
	ID   int                    `json:"id"`
	Rect *physics.KinematicRect `json:"rect"`
}

type DeleteData struct {
	ID int `json:"id"`
}

func (gl *PlatformerGL) NewLevelMessage() *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: LevelData{
			StaticRects:    gl.level.physEng.GetStaticRects(),
			KinematicRects: gl.level.physEng.GetKinematicRects(),
		},
	}
}

func (gl *PlatformerGL) NewDeleteMessage(rectID int) *message.JSON {
	return &message.JSON{
		Type: "delete",
		Body: DeleteData{
			ID: rectID,
		},
	}
}

func (gl *PlatformerGL) NewCreateMessage(rectID int, rect *physics.KinematicRect) *message.JSON {
	return &message.JSON{
		Type: "create",
		Body: CreateData{
			ID:   rectID,
			Rect: rect,
		},
	}
}

func (gl *PlatformerGL) NewGameInfoMessage(playerRectID int, tps int, constants Constants) *message.JSON {
	return &message.JSON{
		Type: "gameinfo",
		Body: GameInfo{
			PlayerRectID: playerRectID,
			Constants:    constants,
			TPS:          tps,
		},
	}
}
