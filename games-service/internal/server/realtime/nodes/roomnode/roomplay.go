package roomnode

import (
	"github.com/neinBit/ocg-games-service/internal/gamelogic"
	"github.com/neinBit/ocg-games-service/internal/gamelogic/clicker"
	"github.com/neinBit/ocg-games-service/internal/gamelogic/platformer"
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
