package platformer

import (
	"github.com/1001bit/OnlineCanvasGames/pkg/set"
)

func (l *Level) CreatePlayer(userID int, playersLimit int) (rectID, *Player) {
	if playerData, ok := l.playersData.Get(userID); ok {
		return playerData.rectID, playerData.player
	}

	rectID := l.getFreePlayerRectID(playersLimit)
	player := NewPlayer(rectID)

	l.players[rectID] = player
	l.playersData.Set(userID, NewPlayerData(player, rectID))

	return rectID, player
}

func (l *Level) DeletePlayer(userID int) (rectID, error) {
	playerData, ok := l.playersData.Get(userID)
	if !ok {
		return 0, ErrNoPlayer
	}

	l.playersData.Delete(userID)
	delete(l.players, playerData.rectID)

	return playerData.rectID, nil
}

func (l *Level) getFreePlayerRectID(playersLimit int) rectID {
	occupiedRectIDs := make(set.Set[rectID])

	playersData, rUnlockFunc := l.playersData.GetMapForRead()
	defer rUnlockFunc()

	for _, playerData := range playersData {
		occupiedRectIDs.Insert(playerData.rectID)
	}

	for newID := range playersLimit {
		if !occupiedRectIDs.Has(rectID(newID)) {
			return rectID(newID)
		}
	}

	return -1
}