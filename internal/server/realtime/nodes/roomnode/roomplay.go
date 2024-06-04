package roomnode

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay/clicker"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/roomplay/platformer"
)

func NewRoomPlayByID(gameID int) roomplay.RoomPlay {
	switch gameID {
	case 1:
		return clicker.NewClickerRP()
	case 2:
		return platformer.NewPlatformerRP()
	default:
		return nil
	}
}
