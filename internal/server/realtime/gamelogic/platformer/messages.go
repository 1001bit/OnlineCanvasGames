package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

type GameInfo struct {
	PlayerRectID int `json:"rectID"`
}

type LevelData struct {
	StaticRects    map[int]*physics.Rect          `json:"static"`
	KinematicRects map[int]*physics.KinematicRect `json:"kinematic"`
}

type CreateData struct {
	ID   int                    `json:"id"`
	Rect *physics.KinematicRect `json:"rect"`
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

func (gl *PlatformerGL) NewDeltasMessage(deltas map[int]*physics.KinematicRect) *message.JSON {
	return &message.JSON{
		Type: "deltas",
		Body: deltas,
	}
}

func (gl *PlatformerGL) NewDeleteMessage(rectID int) *message.JSON {
	return &message.JSON{
		Type: "delete",
		Body: rectID,
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

func (gl *PlatformerGL) NewGameInfoMessage(playerRectID int) *message.JSON {
	return &message.JSON{
		Type: "gameinfo",
		Body: GameInfo{
			PlayerRectID: playerRectID,
		},
	}
}
