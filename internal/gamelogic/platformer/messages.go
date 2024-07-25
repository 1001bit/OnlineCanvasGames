package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/mathobjects"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

// Level
type LevelData struct {
	Blocks  map[int]*Block  `json:"blocks"`
	Players map[int]*Player `json:"players"`

	Config LevelConfig `json:"config"`

	TPS       float64 `json:"tps"`
	ClientTPS float64 `json:"clientTps"`

	PlayerRectID int `json:"playerRectId"`
}

func NewLevelMessage(level *Level, playerRectID int) *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: LevelData{
			Blocks:  level.blocks,
			Players: level.players,

			Config: level.config,

			TPS:       level.tps,
			ClientTPS: level.clientTPS,

			PlayerRectID: playerRectID,
		},
	}
}

// Level Update
type LevelUpdateDate struct {
	MovedPlayers map[int]mathobjects.Vector2[float64] `json:"movedPlayers"`
}

func NewLevelUpdateMessage(movedPlayers map[int]mathobjects.Vector2[float64]) *message.JSON {
	return &message.JSON{
		Type: "levelUpdate",
		Body: LevelUpdateDate{
			MovedPlayers: movedPlayers,
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
