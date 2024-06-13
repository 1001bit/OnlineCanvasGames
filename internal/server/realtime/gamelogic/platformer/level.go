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

	return level
}

func (l *Level) CreatePlayer(userID int, playersLimit int) int {
	const (
		applyGravity    = false
		applyCollisions = true
		applyFriction   = true
	)

	if rectID, ok := l.playersRects[userID]; ok {
		return rectID
	}

	rectID := l.getFreePlayerRectID(playersLimit)

	inner := physics.MakeRect(100*float64(rectID), 100, 100, 100)
	kinRect := physics.NewKinematicRect(inner, applyGravity, applyCollisions, applyFriction)

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

func (l *Level) getFreePlayerRectID(playersLimit int) int {
	occupiedIDs := make([]bool, playersLimit)
	for _, v := range l.playersRects {
		occupiedIDs[v] = true
	}

	for id, occupied := range occupiedIDs {
		if !occupied {
			return id
		}
	}
	return -1
}
