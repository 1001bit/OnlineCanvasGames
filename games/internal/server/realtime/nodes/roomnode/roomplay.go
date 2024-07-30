package roomnode

import (
	"github.com/neinBit/ocg-games-service/internal/gamelogic"
	"github.com/neinBit/ocg-games-service/internal/gamelogic/clicker"
	"github.com/neinBit/ocg-games-service/internal/gamelogic/platformer"
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
