package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

type Level struct {
	physEnv     physics.Environment
	publicRects map[int]*physics.Rect
}

func NewPlatformerLevel() *Level {
	level := &Level{
		physEnv:     *physics.NewEnvironment(),
		publicRects: make(map[int]*physics.Rect),
	}

	staticRect := physics.NewRect(100, 100, 100, 100)
	level.InsertRect(staticRect, 1, true)

	return level
}

func (l *Level) InsertRect(r *physics.Rect, id int, public bool) {
	l.physEnv.InsertRect(r, id)
	if public {
		l.publicRects[id] = r
	}
}

func (l *Level) GetPublicRects() map[int]*physics.Rect {
	return l.publicRects
}
