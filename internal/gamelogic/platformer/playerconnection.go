package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/pkg/set"
)

func (l *Level) CreatePlayer(userID int, playersLimit int) (int, *Player) {
	if rectID, ok := l.userRectIDs[userID]; ok {
		return rectID, l.players[rectID]
	}

	rectID := l.getFreePlayerRectID(playersLimit)
	l.userRectIDs[userID] = rectID

	player := NewPlayer(rectID)
	l.players[rectID] = player

	return rectID, player
}

func (l *Level) DeletePlayer(userID int) (int, error) {
	rectID, ok := l.userRectIDs[userID]
	if !ok {
		return 0, ErrNoPlayer
	}

	delete(l.userRectIDs, userID)
	delete(l.players, rectID)

	return rectID, nil
}

func (l *Level) getFreePlayerRectID(playersLimit int) int {
	occupiedRectIDs := make(set.Set[int])

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
