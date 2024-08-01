package platformer

import (
	"github.com/1001bit/ocg-games-service/internal/mathobjects"
	"github.com/1001bit/ocg-games-service/internal/server/message"
)

// Level
type LevelData struct {
	Blocks  map[rectID]*Block  `json:"blocks"`
	Players map[rectID]*Player `json:"players"`

	Config LevelConfig `json:"config"`

	ServerTPS float64 `json:"serverTps"`
	ClientTPS float64 `json:"clientTps"`

	PlayerRectID rectID `json:"playerRectId"`
}

func NewLevelMessage(level *Level, playerRectID rectID) *message.JSON {
	return &message.JSON{
		Type: "level",
		Body: LevelData{
			Blocks:  level.blocks,
			Players: level.players,

			Config: level.config,

			ServerTPS: level.serverTPS,
			ClientTPS: level.clientTPS,

			PlayerRectID: playerRectID,
		},
	}
}

// Level Update
type LevelUpdateDate struct {
	// rectID[position]
	Players map[rectID]mathobjects.Vector2[float64] `json:"players"`

	DoCorrect bool `json:"correct"`
}

func NewLevelUpdateMessage(sentPlayers map[rectID]mathobjects.Vector2[float64], doCorrect bool) *message.JSON {
	return &message.JSON{
		Type: "levelUpdate",
		Body: LevelUpdateDate{
			Players: sentPlayers,

			DoCorrect: doCorrect,
		},
	}
}

// Connect
type ConnectData struct {
	RectID rectID  `json:"rectId"`
	Player *Player `json:"rect"`
}

func NewConnectMessage(rectID rectID, player *Player) *message.JSON {
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
	RectID rectID `json:"rectId"`
}

func NewDisconnectMessage(rectID rectID) *message.JSON {
	return &message.JSON{
		Type: "disconnect",
		Body: DisconnectData{
			RectID: rectID,
		},
	}
}
