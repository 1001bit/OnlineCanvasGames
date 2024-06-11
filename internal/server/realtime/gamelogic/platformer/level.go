package platformer

import (
	"errors"

	"github.com/1001bit/OnlineCanvasGames/internal/physics"
)

var ErrNoPlayer = errors.New("no such player")

const (
	friction = 0.001
	gForce   = 1
)

type Level struct {
	physEnv physics.Environment

	playersRects map[int]int
}

func NewPlatformerLevel() *Level {
	level := &Level{
		physEnv: *physics.NewEnvironment(friction, gForce),

		playersRects: make(map[int]int),
	}

	return level
}

func (l *Level) CreatePlayer(userID int) int {
	if rectID, ok := l.playersRects[userID]; ok {
		return rectID
	}

	rectID := len(l.playersRects)

	inner := physics.MakeRect(100*float64(rectID), 100, 100, 100)
	kinRect := physics.NewKinematicRect(inner, false, false, true)

	l.physEnv.InsertKinematicRect(kinRect, rectID)

	l.playersRects[userID] = rectID

	return rectID
}

func (l *Level) DeletePlayer(userID int) (int, error) {
	rectID, ok := l.playersRects[userID]
	if !ok {
		return 0, ErrNoPlayer
	}

	delete(l.playersRects, userID)
	l.physEnv.DeleteRect(rectID)

	return rectID, nil
}
