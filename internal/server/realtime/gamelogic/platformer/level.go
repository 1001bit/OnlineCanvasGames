package platformer

import (
	"errors"

	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

var ErrNoPlayer = errors.New("no such player")

type Level struct {
	physEnv physics.Environment

	playersRects map[int]int
}

func NewPlatformerLevel() *Level {
	level := &Level{
		physEnv: *physics.NewEnvironment(),

		playersRects: make(map[int]int),
	}

	block := physics.MakeRect(0, 500, 1000, 100, true)
	level.physEnv.InsertRect(&block, 10)

	return level
}
