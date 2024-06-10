package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

const (
	friction = 0.9
	gForce   = 1
)

type Level struct {
	physEnv     physics.Environment
	publicRects map[int]*physics.Rect

	playersRects map[int]int
}

func NewPlatformerLevel() *Level {
	level := &Level{
		physEnv:     *physics.NewEnvironment(friction, gForce),
		publicRects: make(map[int]*physics.Rect),

		playersRects: make(map[int]int),
	}

	return level
}

func (l *Level) InsertRect(r *physics.Rect, id int, public bool) {
	l.physEnv.InsertRect(r, id)
	if public {
		l.publicRects[id] = r
	}
}

func (l *Level) InsertKinematicRect(kr *physics.KinematicRect, id int, public bool) {
	l.physEnv.InsertKinematicRect(kr, id)
	if public {
		l.publicRects[id] = kr.GetRect()
	}
}

func (l *Level) DeleteRect(id int) {
	l.physEnv.DeleteRect(id)
	delete(l.publicRects, id)
}

func (l *Level) CreatePlayer(userID int) int {
	if rectID, ok := l.playersRects[userID]; ok {
		return rectID
	}

	rectID := len(l.playersRects)

	inner := physics.NewRect(100*float64(rectID), 100, 100, 100)
	kinRect := physics.NewKinematicRect(inner, false, false, true)

	l.InsertKinematicRect(kinRect, rectID, true)

	l.playersRects[userID] = rectID

	return rectID
}

func (l *Level) DeletePlayer(userID int) {
	rectID, ok := l.playersRects[userID]
	if !ok {
		return
	}

	delete(l.playersRects, userID)
	l.DeleteRect(rectID)
}

func (l *Level) GetPublicRects() map[int]*physics.Rect {
	return l.publicRects
}
