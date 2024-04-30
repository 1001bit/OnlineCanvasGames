package gamenode

import gamemodel "github.com/1001bit/OnlineCanvasGames/internal/model/game"

// ask gameNode to update gameNode.roomsJSON
func (gameNode *GameNode) RequestUpdatingRoomsJSON() {
	gameNode.roomsJSONUpdateChan <- struct{}{}
}

func (gameNode *GameNode) GetGame() gamemodel.Game {
	return gameNode.game
}
