package roomnode

import (
	"github.com/1001bit/onlinecanvasgames/services/games/internal/gamelogic"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/gamelogic/clicker"
	"github.com/1001bit/onlinecanvasgames/services/games/internal/gamelogic/platformer"
)

func NewGameLogicByTitle(title string) gamelogic.GameLogic {
	switch title {
	case "clicker":
		return clicker.NewClickerGL()
	case "platformer":
		return platformer.NewPlatformerGL()
	default:
		return nil
	}
}
