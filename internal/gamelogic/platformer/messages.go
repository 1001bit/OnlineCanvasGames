package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

// GameInfo
type GameInfo struct {
	TPS int `json:"tps"`
}

func NewGameInfoMessage(gl *PlatformerGL) *message.JSON {
	return &message.JSON{
		Type: "gameinfo",
		Body: GameInfo{
			TPS: gl.tps,
		},
	}
}

// Level
type LevelData struct {
	Blocks  map[int]*Block  `json:"blocks"`
	Players map[int]*Player `json:"players"`

	Config LevelConfig `json:"config"`

	PlayerRectID int `json:"playerRectId"`
}

func NewLevelMessage(level *Level, userID int) *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: LevelData{
			Blocks:  level.blocks,
			Players: level.players,

			Config: level.config,

			PlayerRectID: level.userRectIDs[userID],
		},
	}
}

// Players movement
type PlayerMovementData struct {
	PlayersMoved map[int]mathobjects.Vector2[float64] `json:"playersMoved"`
}

func NewPlayerMovementMessage(movedPlayers map[int]mathobjects.Vector2[float64]) *message.JSON {
	return &message.JSON{
		Type: "playerMovement",
		Body: PlayerMovementData{
			PlayersMoved: movedPlayers,
		},
	}
}

// Connect
type ConnectData struct {
	RectID int     `json:"rectId"`
	Player *Player `json:"rect"`
}

func NewConnectMessage(rectID int, player *Player) *message.JSON {
	return &message.JSON{
		Type: "connect",
		Body: ConnectData{
			RectID: rectID,
			Player: player,
		},
	}
}

// Disconnect
type DisconnectData struct {
	RectID int `json:"rectId"`
}

func NewDisconnectMessage(rectID int) *message.JSON {
	return &message.JSON{
		Type: "disconnect",
		Body: DisconnectData{
			RectID: rectID,
		},
	}
}
