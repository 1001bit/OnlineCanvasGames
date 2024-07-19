package platformer

import (
	"errors"

	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

var ErrNoPlayer = errors.New("no such player")

type Level struct {
	physEng physics.Engine

	playersRects map[int]int
}

func NewPlatformerLevel() *Level {
	level := &Level{
		physEng: *physics.NewEngine(),

		playersRects: make(map[int]int),
	}

	block := physics.MakePhysicalRect(0, 500, 1000, 100, true)
	level.physEng.InsertRect(&block, 10)

	return level
}
