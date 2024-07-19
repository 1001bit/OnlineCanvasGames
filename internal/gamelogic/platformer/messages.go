package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

// GameInfo
type GameInfo struct {
	PlayerRectID int       `json:"rectID"`
	Constants    Constants `json:"constants"`
	TPS          int       `json:"tps"`
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

// Level
type LevelData struct {
	StaticRects    map[int]*physics.PhysicalRect  `json:"static"`
	KinematicRects map[int]*physics.KinematicRect `json:"kinematic"`
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

// Update
type UpdateData struct {
	RectsMoved map[int]mathobjects.Vector2[float64] `json:"rectsMoved"`
}

func (gl *PlatformerGL) NewUpdateMessage(movedRects map[int]mathobjects.Vector2[float64]) *message.JSON {
	return &message.JSON{
		Type: "update",
		Body: UpdateData{
			RectsMoved: movedRects,
		},
	}
}

// Create
type CreateData struct {
	ID   int                    `json:"id"`
	Rect *physics.KinematicRect `json:"rect"`
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

// Delete
type DeleteData struct {
	ID int `json:"id"`
}

func (gl *PlatformerGL) NewDeleteMessage(rectID int) *message.JSON {
	return &message.JSON{
		Type: "delete",
		Body: DeleteData{
			ID: rectID,
		},
	}
}
