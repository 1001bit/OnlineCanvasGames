package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/physics"

func (l *Level) CreatePlayer(userID int, playersLimit int) int {
	const (
		applyGravity    = true
		applyCollisions = true
		applyFriction   = true
	)

	if rectID, ok := l.playersRects[userID]; ok {
		return rectID
	}

	rectID := l.getFreePlayerRectID(playersLimit)

	inner := physics.MakePhysicalRect(100*float64(rectID), 100, 100, 100, true)
	kinRect := physics.NewKinematicRect(inner, applyGravity, applyFriction)

	l.physEng.InsertKinematicRect(kinRect, rectID)

	l.playersRects[userID] = rectID

	return rectID
}

func (l *Level) DeletePlayer(userID int) (int, error) {
	rectID, ok := l.playersRects[userID]
	if !ok {
		return 0, ErrNoPlayer
	}

	delete(l.playersRects, userID)
	l.physEng.DeleteRect(rectID)

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
