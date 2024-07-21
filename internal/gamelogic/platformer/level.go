package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

type Level struct {
	physEng physics.Engine

	// [userID]rectID
	playersRects map[int]int
}

func NewPlatformerLevel() *Level {
	level := &Level{
		physEng: *physics.NewEngine(),

		playersRects: make(map[int]int),
	}

	block := physics.NewPhysicalRect(0, 500, 1000, 100, true)
	level.physEng.InsertStaticRect(block, 10)

	return level
}
