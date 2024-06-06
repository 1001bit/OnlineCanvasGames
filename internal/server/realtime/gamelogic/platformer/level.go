package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/server/realtime/gamelogic"

func NewPlatformerLevel() *gamelogic.Level {
	rects := make(map[int]*gamelogic.Rect)

	static := gamelogic.NewRect(100, 100, 100, 100, false, true)
	rects[1] = static

	return gamelogic.NewLevel(rects)
}
