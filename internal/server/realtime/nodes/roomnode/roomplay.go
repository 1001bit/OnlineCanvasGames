package roomnode

import (
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic/clicker"
	"github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic/platformer"
)

func NewGameLogicByID(gameID int) gamelogic.GameLogic {
	switch gameID {
	case 1:
		return clicker.NewClickerGL()
	case 2:
		return platformer.NewPlatformerGL()
	default:
		return nil
	}
}
