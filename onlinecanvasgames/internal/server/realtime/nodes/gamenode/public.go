package gamenode

import (
	gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"
	"github.com/1001bit/OnlineCanvasGames/internal/server/message"
)

// ask gameNode to update gameNode.roomsJSON
func (gameNode *GameNode) RequestUpdatingRoomsJSON() {
	select {
	case gameNode.roomsJSONUpdateChan <- struct{}{}:
		// Send request to update roomsJSON
	case <-gameNode.Flow.Done():
		// gamenode is done
	}

}

// write a message to every client
func (gameNode *GameNode) GlobalWriteMessage(msg *message.JSON) {
	childrenSet, rUnlockFunc := gameNode.Clients.ChildrenSet.GetSetForRead()
	defer rUnlockFunc()

	for client := range childrenSet {
		go client.WriteMessage(msg)
	}
}

func (gameNode *GameNode) GetGame() gamemodel.Game {
	return gameNode.game
}
