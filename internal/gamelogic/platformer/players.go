package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/internal/physics"
	"github.com/1001bit/OnlineCanvasGames/pkg/set"
)

func (l *Level) CreatePlayer(userID int, playersLimit int) (int, *physics.KinematicRect) {
	const (
		applyGravity    = true
		applyCollisions = true
		applyFriction   = true
	)

	if rectID, ok := l.userRectIDs[userID]; ok {
		return rectID, l.playersRects[rectID]
	}

	rectID := l.getFreePlayerRectID(playersLimit)

	inner := physics.MakePhysicalRect(100*float64(rectID), 100, 100, 100, applyCollisions)
	kinRect := physics.NewKinematicRect(inner, physics.FrictionType, physics.GravityType)

	l.playersRects[rectID] = kinRect
	l.userRectIDs[userID] = rectID

	return rectID, kinRect
}

func (l *Level) DeletePlayer(userID int) (int, error) {
	rectID, ok := l.userRectIDs[userID]
	if !ok {
		return 0, ErrNoPlayer
	}

	delete(l.userRectIDs, userID)
	delete(l.playersRects, rectID)

	return rectID, nil
}

func (l *Level) getFreePlayerRectID(playersLimit int) int {
	occupiedRectIDs := set.MakeEmptySet[int]()

	for _, rectID := range l.userRectIDs {
		occupiedRectIDs.Insert(rectID)
	}

	for newID := range playersLimit {
		if !occupiedRectIDs.Has(newID) {
			return newID
		}
	}
	return -1
}
