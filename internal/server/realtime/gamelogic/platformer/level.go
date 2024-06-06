package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"

func NewPlatformerLevel() *gamelogic.Level {
	level := gamelogic.NewLevel()

	staticRect := gamelogic.NewRect(100, 100, 100, 100)
	level.InsertRect(staticRect, 1, true)

	return level
}
